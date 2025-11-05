#![allow(unexpected_cfgs)]
#![allow(deprecated)]

use anchor_lang::prelude::*;

declare_id!("AdUTQjW9iWgWwjsr7n5RjVLjt1VGNtBSviJQtk18ESxQ");

// Anchor accounts always reserve 8 bytes for the discriminator.
// The discriminator is a fixed prefix set by Anchor to identify the account type.
pub const ANCHOR_DISCRIMINATOR_SIZE: usize = 8;

/// Anchor program that stores a user's favorite number, color and hobbies
/// in a PDA (Program Derived Address) account. The PDA is derived from a
/// static seed ("favorites") and the user's public key. The `initialize`
/// instruction creates the PDA if it doesn't exist, or updates it otherwise.
#[program]
pub mod favorite {
    use super::*;

    /// Creates or updates the user's favorites account.
    ///
    /// - `ctx`: Account context, containing the signer (`user`) and the PDA (`favorites`).
    /// - `number`: The user's favorite number.
    /// - `color`: The user's favorite color (string, limited in length by `InitSpace`).
    /// - `hobbies`: A list of the user's hobbies (vector of strings, each bounded in length).
    ///
    /// Behavior:
    /// - If the PDA does not exist yet, it is created via `init_if_needed` with the provided space.
    /// - The PDA is addressed by seeds `[b"favorites", user.key().as_ref()]` and a bump.
    /// - The instruction writes the provided values into the PDA account.
    pub fn initialize(
        ctx: Context<SetFavorite>,
        number: u64,
        color: String,
        hobbies: Vec<String>,
    ) -> Result<()> {
        let _public_id = ctx.accounts.user.key();
        msg!("Greetings from: {:?}", ctx.program_id);
        msg!(
            "User {public_id}'s favorite number is {number}, 
        favorite color is: {color}"
        );
        ctx.accounts.favorites.set_inner(Favorite {
            number,
            color,
            hobbies,
        });
        Ok(())
    }

    pub fn update(
        ctx: Context<SetFavorite>,
        number: u64,
        color: String,
        hobbies: Vec<String>,
    ) -> Result<()> {

        ctx.accounts.favorites.set_inner(Favorite {
            number,
            color,
            hobbies,
        });

        Ok(())
    }
}

/// On-chain account that stores the user's preferences.
///
/// The `InitSpace` derive macro calculates the static space required for this
/// account based on the field types and the `#[max_len]` annotations.
/// It provides an associated constant `INIT_SPACE` that we can use when
/// initializing the account.
#[account]
#[derive(InitSpace)]
pub struct Favorite {
    /// Favorite number (64-bit unsigned integer).
    pub number: u64,

    /// Favorite color. Bounded string with a maximum length of 10 characters.
    #[max_len(10)]
    pub color: String,

    /// Favorite hobbies. A vector with at most 10 strings,
    /// where each string is at most 50 characters.
    #[max_len(10, 50)]
    pub hobbies: Vec<String>,
}

/// Account context for the `initialize` instruction.
///
/// - `user`: The transaction signer and payer.
/// - `favorites`: PDA that stores the `Favorite` data. Created if missing.
/// - `system_program`: Required for account creation and rent payments.
#[derive(Accounts)]
pub struct SetFavorite<'info> {
    #[account(mut)]
    pub user: Signer<'info>,

    #[account(
        init_if_needed,
        payer = user,
        // Allocate enough space for the Anchor discriminator plus the account data.
        // `InitSpace` derive exposes `Favorite::INIT_SPACE` for the data portion.
        space = ANCHOR_DISCRIMINATOR_SIZE + Favorite::INIT_SPACE,
        // PDA seeds: static label + user's public key
        seeds = [b"favorites", user.key().as_ref()],
        // Auto-derive the PDA bump for the given seeds
        bump
    )]
    pub favorites: Account<'info, Favorite>,

    pub system_program: Program<'info, System>,
}
