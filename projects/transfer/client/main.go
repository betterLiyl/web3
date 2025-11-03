package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/mr-tron/base58" 
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/sysprog"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/blocto/solana-go-sdk/types"
)

const lamportsPerSOL = 1_000_000_000

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	sub := os.Args[1]
	switch sub {
	case "balance":
		balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
		addr := balanceCmd.String("address", "", "Account base58 address")
		cluster := balanceCmd.String("cluster", "devnet", "Cluster: devnet|testnet|mainnet|local")
		rpc := balanceCmd.String("rpc", "", "Custom RPC endpoint URL (override)")
		_ = balanceCmd.Parse(os.Args[2:])
		if *addr == "" {
			log.Fatal("missing --address")
		}
		if !isValidBase58Pubkey(*addr) {
			log.Fatal("invalid --address base58")
		}
		if err := runBalance(*addr, normalizeCluster(*cluster), strings.TrimSpace(*rpc)); err != nil {
			log.Fatalf("balance error: %v", err)
		}
	case "transfer":
		transferCmd := flag.NewFlagSet("transfer", flag.ExitOnError)
		fromPriv := transferCmd.String("from", "", "Sender private key (base58)")
		fromFile := transferCmd.String("fromFile", "", "Path to Solana keypair JSON file (id.json)")
		toAddr := transferCmd.String("to", "", "Recipient address (base58)")
		lamports := transferCmd.Uint64("lamports", 0, "Amount in lamports (1 SOL = 1e9)")
		cluster := transferCmd.String("cluster", "devnet", "Cluster: devnet|testnet|mainnet|local")
		rpc := transferCmd.String("rpc", "", "Custom RPC endpoint URL (override)")
		_ = transferCmd.Parse(os.Args[2:])
		if (*fromPriv == "" && *fromFile == "") || *toAddr == "" || *lamports == 0 {
			log.Fatal("missing required flags: --from or --fromFile, --to, --lamports")
		}
		if !isValidBase58Pubkey(*toAddr) {
			log.Fatal("invalid --to base58")
		}
		if err := runTransfer(*fromPriv, *fromFile, *toAddr, *lamports, normalizeCluster(*cluster), strings.TrimSpace(*rpc)); err != nil {
			log.Fatalf("transfer error: %v", err)
		}
	case "airdrop":
		airdropCmd := flag.NewFlagSet("airdrop", flag.ExitOnError)
		toAddr := airdropCmd.String("to", "", "Recipient address (base58)")
		lamports := airdropCmd.Uint64("lamports", lamportsPerSOL, "Amount in lamports (default 1 SOL)")
		cluster := airdropCmd.String("cluster", "local", "Cluster: devnet|testnet|mainnet|local")
		rpc := airdropCmd.String("rpc", "", "Custom RPC endpoint URL (override)")
		_ = airdropCmd.Parse(os.Args[2:])
		if *toAddr == "" || *lamports == 0 {
			log.Fatal("missing required flags: --to, --lamports")
		}
		if !isValidBase58Pubkey(*toAddr) {
			log.Fatal("invalid --to base58")
		}
		if err := runAirdrop(*toAddr, *lamports, normalizeCluster(*cluster), strings.TrimSpace(*rpc)); err != nil {
			log.Fatalf("airdrop error: %v", err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func runBalance(address string, cluster string, rpcOverride string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c := client.NewClient(resolveEndpoint(cluster, rpcOverride))
	bal, err := c.GetBalance(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}
	lam := uint64(bal)
	sol := new(big.Float).Quo(new(big.Float).SetUint64(lam), new(big.Float).SetUint64(lamportsPerSOL))
	solF, _ := sol.Float64()

	out := map[string]any{
		"address":  address,
		"cluster":  cluster,
		"lamports": lam,
		"sol":      solF,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func runTransfer(fromPrivBase58, fromFilePath, toAddrBase58 string, amountLamports uint64, cluster string, rpcOverride string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()
	var from types.Account
	var err error
	if strings.TrimSpace(fromPrivBase58) != "" {
		from, err = types.AccountFromBase58(strings.TrimSpace(fromPrivBase58))
		if err != nil {
			return fmt.Errorf("invalid sender private key: %w", err)
		}
	} else {
		from, err = loadAccountFromFile(strings.TrimSpace(fromFilePath))
		if err != nil {
			return fmt.Errorf("failed to load sender from file: %w", err)
		}
	}
	// validate recipient and parse
	if !isValidBase58Pubkey(toAddrBase58) {
		return errors.New("recipient address invalid")
	}
	to := common.PublicKeyFromString(strings.TrimSpace(toAddrBase58))
	c := client.NewClient(resolveEndpoint(cluster, rpcOverride))
	latest, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest blockhash: %w", err)
	}
	recent := latest.Blockhash

	msg := types.NewMessage(types.NewMessageParam{
		FeePayer:        from.PublicKey,
		RecentBlockhash: recent,
		Instructions: []types.Instruction{
			sysprog.Transfer(sysprog.TransferParam{
				From:   from.PublicKey,
				To:     to,
				Amount: amountLamports,
			}),
		},
	})

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: msg,
		Signers: []types.Account{from},
	})
	if err != nil {
		return fmt.Errorf("failed to build transaction: %w", err)
	}

	txhash, err := c.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	out := map[string]any{
		"cluster":   cluster,
		"blockhash": recent,
		"txhash":    txhash,
		"amount":    amountLamports,
		"from":      from.PublicKey.ToBase58(),
		"to":        to.ToBase58(),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func runAirdrop(toAddrBase58 string, lamports uint64, cluster string, rpcOverride string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	to := strings.TrimSpace(toAddrBase58)
	c := client.NewClient(resolveEndpoint(cluster, rpcOverride))
	txhash, err := c.RequestAirdrop(ctx, to, lamports)
	if err != nil {
		return fmt.Errorf("failed to request airdrop: %w", err)
	}
	out := map[string]any{
		"cluster":  cluster,
		"txhash":   txhash,
		"lamports": lamports,
		"to":       to,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func endpointFor(cluster string) string {
	switch cluster {
	case "mainnet", "mainnet-beta":
		return rpc.MainnetRPCEndpoint
	case "testnet":
		return rpc.TestnetRPCEndpoint
	case "local":
		return "http://127.0.0.1:8899"
	default:
		return rpc.DevnetRPCEndpoint
	}
}

func resolveEndpoint(cluster, rpcOverride string) string {
	if strings.TrimSpace(rpcOverride) != "" {
		return strings.TrimSpace(rpcOverride)
	}
	return endpointFor(cluster)
}

func normalizeCluster(c string) string {
	c = strings.ToLower(strings.TrimSpace(c))
	switch c {
	case "mainnet", "mainnet-beta":
		return "mainnet"
	case "testnet":
		return "testnet"
	case "local":
		return "local"
	default:
		return "devnet"
	}
}

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

func printUsage() {
	fmt.Println(`Usage:
  Balance:
    go run main.go balance --address <base58> [--cluster devnet|testnet|mainnet|local] [--rpc <url>]

  Transfer SOL:
    go run main.go transfer (--from <privateKeyBase58> | --fromFile ~/.config/solana/id.json) --to <addressBase58> --lamports <amount> [--cluster devnet|testnet|mainnet|local] [--rpc <url>]

  Airdrop (devnet/local only):
    go run main.go airdrop --to <addressBase58> [--lamports 1000000000] [--cluster local|devnet] [--rpc <url>]`)
}

func isValidBase58Pubkey(s string) bool {
	b, err := base58.Decode(strings.TrimSpace(s))
	return err == nil && len(b) == 32
}