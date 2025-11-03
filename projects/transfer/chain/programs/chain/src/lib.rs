#![allow(unexpected_cfgs)]
#![allow(deprecated)]

use anchor_lang::prelude::*;
use anchor_lang::system_program::{self, System};

// 本程序演示：
// 1) 在 Solana 上通过 PDA 派生创建一个系统钱包账户（SystemAccount）
// 2) 查询任意系统钱包账户的 SOL 余额（lamports）
//
// 说明：
// - 钱包账户是由系统程序（System Program）拥有的“普通账户”。
// - 本程序通过 PDA（Program Derived Address）派生出一个地址，并创建为系统账户。
//   这样无需外部私钥签名，程序即可控制它（通过 seeds + bump 进行“签名”）。
// - 合约无法返回数据，只能通过日志（msg!）或事件（emit!）向客户端传递信息。

declare_id!("A41gXaRcvZDSFEf2vLg1wjxKwi3ybbT3f2yvd6ZhBYer");

#[allow(deprecated)]
#[program]
pub mod chain {
    use super::*;

    /// 创建一个系统钱包账户（PDA），并可选择转入初始 SOL。
    ///
    /// 该账户是系统程序（System Program）拥有的普通钱包账户，
    /// 地址由本合约的种子派生（PDA），因此只有本程序能“签名”控制它。
    ///
    /// 参数:
    /// - `seed`：自定义字符串，用于参与 PDA 派生，保证地址可预测
    /// - `initial_lamports`：创建后从 `payer` 向新钱包转入的 SOL（单位 lamports）
    pub fn create_wallet(
        ctx: Context<CreateWallet>,
        _seed: String,
        initial_lamports: u64,
    ) -> Result<()> {
        // 使用 PDA 签名手动创建一个系统账户（space=0，owner=System），
        // 并由 payer 支付 `initial_lamports` 作为初始余额。
        let cpi_accounts = system_program::CreateAccount {
            from: ctx.accounts.payer.to_account_info(),
            to: ctx.accounts.wallet.to_account_info(),
        };

        // 通过 seeds + bump 进行程序签名，以允许创建 PDA 账户。
        let wallet_seeds: &[&[u8]] = &[
            b"wallet",
            ctx.accounts.payer.key.as_ref(),
            &[ctx.bumps.wallet],
        ];
        // 注意：不能直接传 `&[wallet_seeds]`，那是临时值，
        // 会在 `CpiContext` 持有引用后被释放，从而触发借用错误。
        let signer_seeds: &[&[&[u8]]] = &[wallet_seeds];
        let cpi_ctx = CpiContext::new_with_signer(
            ctx.accounts.system_program.to_account_info(),
            cpi_accounts,
            signer_seeds,
        );
        // 创建系统账户（数据空间=0，所有者=System Program）。
        system_program::create_account(cpi_ctx, initial_lamports, 0, &System::id())?;

        // 输出日志，方便在交易日志中查看。
        let balance = ctx.accounts.wallet.to_account_info().lamports();
        msg!(
            "Wallet created at: {:?}, current balance: {} lamports",
            ctx.accounts.wallet.key(),
            balance
        );

        Ok(())
    }

    /// 查询任意系统钱包账户的 SOL 余额（以 lamports 计）。
    ///
    /// 说明：
    /// - 合约无法“返回”值，只能通过日志或事件输出；
    /// - 这里同时通过 `msg!` 打印，并发出 `BalanceEvent` 事件以便客户端解析。
    pub fn get_balance(ctx: Context<GetBalance>) -> Result<()> {
        let lamports = ctx.accounts.wallet.to_account_info().lamports();
        msg!(
            "Balance of {:?}: {} lamports",
            ctx.accounts.wallet.key(),
            lamports
        );
        emit!(BalanceEvent {
            wallet: ctx.accounts.wallet.key(),
            lamports
        });
        Ok(())
    }
}

/// 创建系统钱包账户时使用的账户集合。
///
/// 关键点：
/// - `payer`：实际出资者，支付新账户的租金以及可选的初始存入金额；必须是签名者
/// - `wallet`：要创建的系统账户（SystemAccount）。使用 PDA 派生，确保地址可预测；`space=0` 因为系统账户无额外数据
/// - `system_program`：系统程序，用于 CPI 创建账户与转账
#[derive(Accounts)]
pub struct CreateWallet<'info> {
    /// 出资者（签名者），为新账户支付租金与初始转账
    #[account(mut)]
    pub payer: Signer<'info>,

    /// 系统钱包账户（PDA）
    /// - 使用 `seeds`+`bump` 派生地址；
    /// - 这里改为使用 `UncheckedAccount` 并在指令中手动 CPI 创建，
    ///   以规避某些环境下 `#[derive(Accounts)]` 的宏展开问题。
    /// CHECK: 该账户地址由本程序使用 `seeds = [b"wallet", payer.key()]` 与 `bump`
    /// 派生（PDA），并在指令中通过对系统程序的 CPI 创建为系统账户（owner=System,
    /// space=0）。我们仅读取其 `lamports`，不依赖任何自定义数据结构，因此无需
    /// 通过类型系统进一步校验。
    #[account(
        mut,
        seeds = [b"wallet", payer.key().as_ref()],
        bump
    )]
    pub wallet: UncheckedAccount<'info>,

    /// 系统程序
    pub system_program: Program<'info, System>,
}

/// 查询余额时使用的账户集合。
///
/// - 这里只需要被查询的系统账户即可（无需签名）
/// - 为了泛用性，不强制该账户必须是本程序生成的 PDA，只要是系统账户即可
#[derive(Accounts)]
pub struct GetBalance<'info> {
    /// 要查询余额的系统账户
    pub wallet: SystemAccount<'info>,
}

/// 用于在链上事件中输出余额信息，方便客户端解析。
#[event]
pub struct BalanceEvent {
    /// 被查询的钱包地址
    pub wallet: Pubkey,
    /// 余额（单位：lamports；1 SOL = 1_000_000_000 lamports）
    pub lamports: u64,
}
