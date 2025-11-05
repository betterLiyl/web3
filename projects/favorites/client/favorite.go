package main

import (
    "context"
    "crypto/sha256"
    "encoding/binary"
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/blocto/solana-go-sdk/client"
    "github.com/blocto/solana-go-sdk/common"
    "github.com/blocto/solana-go-sdk/types"
)

// 已部署 Anchor 程序的 Program ID（与 declare_id! 一致）
const favoriteProgramID = "AdUTQjW9iWgWwjsr7n5RjVLjt1VGNtBSviJQtk18ESxQ"

// 默认使用 Devnet，如需切换可改为本地或主网 RPC
const defaultEndpoint = "https://api.devnet.solana.com"

func main() {
    // 加载签名者：~/.config/solana/id.json
    signer, err := loadAccountFromFile(os.ExpandEnv("$HOME/.config/solana/id.json"))
    if err != nil {
        fmt.Printf("failed to load signer: %v\n", err)
        return
    }

    c := client.NewClient(defaultEndpoint)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // 预检余额，devnet 可自动空投确保交易与账户创建费用
    if err := ensureAirdropIfLow(ctx, c, signer.PublicKey.ToBase58(), 500_000_000 /* 0.5 SOL */); err != nil {
        fmt.Printf("airdrop check failed: %v\n", err)
        // 不中断主流程
    }

    // 获取最新区块哈希
    latest, err := c.GetLatestBlockhash(ctx)
    if err != nil {
        fmt.Printf("failed to get latest blockhash: %v\n", err)
        return
    }
    recent := latest.Blockhash

    // 派生 PDA：seeds = ["favorites", user]
    programID := common.PublicKeyFromString(favoriteProgramID)
    // 注意：blocto SDK 的 FindProgramAddress 接收 [][]byte 作为 seeds
    // PublicKey.Bytes() 提供原始 32 字节公钥
    favoritesPDA, _, err := common.FindProgramAddress(
        [][]byte{[]byte("favorites"), signer.PublicKey.Bytes()},
        programID,
    )
    if err != nil {
        fmt.Printf("failed to find PDA: %v\n", err)
        return
    }

    // 构造 initialize 指令数据（Anchor: 8 字节 discriminator + Borsh 编码参数）
    ixData := encodeInitialize(42, "blue", []string{"reading", "coding"})

    // 构建指令账户列表：顺序需与 SetFavorite 定义一致
    metas := []types.AccountMeta{
        {PubKey: signer.PublicKey, IsSigner: true, IsWritable: true},
        {PubKey: favoritesPDA, IsSigner: false, IsWritable: true},
        // system_program
        {PubKey: common.PublicKeyFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: false},
    }

    ix := types.Instruction{
        ProgramID: programID,
        Accounts:  metas,
        Data:      ixData,
    }

    // 构建并签名交易
    msg := types.NewMessage(types.NewMessageParam{
        FeePayer:        signer.PublicKey,
        RecentBlockhash: recent,
        Instructions:    []types.Instruction{ix},
    })

    tx, err := types.NewTransaction(types.NewTransactionParam{
        Message: msg,
        Signers: []types.Account{signer},
    })
    if err != nil {
        fmt.Printf("failed to build tx: %v\n", err)
        return
    }

    sig, err := c.SendTransaction(ctx, tx)
    if err != nil {
        fmt.Printf("failed to send tx: %v\n", err)
        return
    }
    fmt.Println("initialize tx signature:", sig)
}

// encodeInitialize 生成 Anchor 指令数据：
// [8字节 discriminator("global:initialize")] [u64 number] [borsh string color] [borsh Vec<String> hobbies]
func encodeInitialize(number uint64, color string, hobbies []string) []byte {
    disc := anchorDiscriminator("initialize")
    // u64（LE）
    num := make([]byte, 8)
    binary.LittleEndian.PutUint64(num, number)
    // color（Borsh string）
    colorB := borshEncodeString(color)
    // hobbies（Borsh Vec<String>）
    hobbiesB := borshEncodeStringVec(hobbies)
    out := make([]byte, 0, len(disc)+len(num)+len(colorB)+len(hobbiesB))
    out = append(out, disc...)
    out = append(out, num...)
    out = append(out, colorB...)
    out = append(out, hobbiesB...)
    return out
}

// anchorDiscriminator 计算 8 字节 SIGHASH("global:<name>")
func anchorDiscriminator(name string) []byte {
    s := "global:" + strings.TrimSpace(name)
    h := sha256.Sum256([]byte(s))
    // 取前 8 字节
    return h[:8]
}

// Borsh: string = u32(len) + utf8 bytes
func borshEncodeString(s string) []byte {
    bs := []byte(s)
    out := make([]byte, 4+len(bs))
    binary.LittleEndian.PutUint32(out[:4], uint32(len(bs)))
    copy(out[4:], bs)
    return out
}

// Borsh: Vec<String> = u32(len) + each string
func borshEncodeStringVec(v []string) []byte {
    out := make([]byte, 4)
    binary.LittleEndian.PutUint32(out[:4], uint32(len(v)))
    for _, s := range v {
        sb := borshEncodeString(s)
        out = append(out, sb...)
    }
    return out
}

// ensureAirdropIfLow: devnet 余额低于阈值时尝试空投
func ensureAirdropIfLow(ctx context.Context, c *client.Client, addr string, min uint64) error {
    bal, err := c.GetBalance(ctx, addr)
    if err != nil {
        return err
    }
    if uint64(bal) >= min {
        return nil
    }
    // 请求 1 SOL 空投
    _, err = c.RequestAirdrop(ctx, addr, 1_000_000_000)
    return err
}

// 复用 transfer 客户端的本地密钥读取逻辑
func loadAccountFromFile(path string) (types.Account, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return types.Account{}, err
    }
    var ints []int
    if err := json.Unmarshal(data, &ints); err != nil {
        return types.Account{}, err
    }
    b := make([]byte, len(ints))
    for i, v := range ints {
        b[i] = byte(v)
    }
    return types.AccountFromBytes(b)
}