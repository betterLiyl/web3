#![allow(unexpected_cfgs)]
#![allow(deprecated)]

use anchor_lang::{prelude::*, solana_program::message};

declare_id!("31Tq6cGFa1CU8JaU51snTvKaXaKqWP3M3dFBWNXeJqYj");

#[program]
pub mod voting {
    use super::*;

    pub fn initialize_poll(
        ctx: Context<InitializePoll>, 
        _poll_id: u64,
        start: u64,
        end: u64,
        name: String,
        desc: String) -> Result<()> {
        let poll = &mut ctx.accounts.poll_account;
        poll.poll_name = name;
        poll.poll_desc = desc;
        poll.poll_vote_start = start;
        poll.poll_vote_end = end;
        poll.poll_vote_index = 0;
        Ok(())
    }

    pub fn initialize_candidate(ctx: Context<InitializeCandidate>,_poll_id: u64,candidate:String) -> Result<()>{
        
        ctx.accounts.candidate_account.candidate_name = candidate;

        ctx.accounts.poll_account.poll_vote_index += 1;
        
        Ok(())
    }


    pub fn vote(ctx: Context<Vote>,_poll_id: u64,_candidate:String) -> Result<()>{
        let candidate = &mut ctx.accounts.candidate_account;
        let current_time = Clock::get()?.unix_timestamp;
        if current_time > (ctx.accounts.poll_account.poll_vote_end as i64) {
            return Err(VotingErrorCode::VotingEnded.into());
        }
        if current_time < (ctx.accounts.poll_account.poll_vote_start as i64) {
            return Err(VotingErrorCode::VotingNotStarted.into());
        }

        candidate.candidate_votes += 1;

        Ok(())
    }


}

#[derive(Accounts)]
#[instruction(poll_id: u64,candidate:String)]
pub struct Vote<'info>{

    #[account(mut)]
    pub signer:Signer<'info>,

    #[account(
        mut,
        seeds = [b"poll", poll_id.to_le_bytes().as_ref()],
        bump
    )]
    pub poll_account: Account<'info, Poll>,

    #[account(
        mut,
        seeds = [poll_id.to_le_bytes().as_ref(),candidate.as_ref()],
        bump
    )]
    pub candidate_account: Account<'info,CandidateAccount>,


}


#[derive(Accounts)]
#[instruction(poll_id: u64)]
pub struct InitializePoll<'info> {

    #[account(mut)]
    pub signer: Signer<'info>,

    #[account(
        init_if_needed,
         payer = signer, 
         space = 8 + Poll::INIT_SPACE,
         seeds = [b"poll", poll_id.to_le_bytes().as_ref()],
         bump,
        )]
    pub poll_account: Account<'info, Poll>,

    pub system_program: Program<'info, System>,
}

#[account]
#[derive(InitSpace)]
pub struct Poll {
    #[max_len(10)]
    pub poll_name: String,

    #[max_len(100)]
    pub poll_desc: String,

    pub poll_vote_start: u64,

    pub poll_vote_end: u64,

    pub poll_vote_index: u64,
}
#[derive(Accounts)]
#[instruction(poll_id:u64,candidate:String)]
pub struct InitializeCandidate<'info>{
    
    #[account(mut)]
    pub signer: Signer<'info>,


    pub poll_account: Account<'info, Poll>,

    #[account(
        init_if_needed,
        payer = signer,
        space = 8 + CandidateAccount::INIT_SPACE,
        seeds = [poll_id.to_le_bytes().as_ref(),candidate.as_ref()],
        bump
    )]
    pub candidate_account: Account<'info,CandidateAccount>,


    pub system_program: Program<'info,System>



}




#[account]
#[derive(InitSpace)]
pub struct CandidateAccount{

    #[max_len(10)]
    pub candidate_name: String,

    pub candidate_votes: u64
}

#[error_code]
pub enum VotingErrorCode {
    
    #[msg("Voting has not started yet")]
    VotingNotStarted,

    #[msg("Voting has ended")]
    VotingEnded
}