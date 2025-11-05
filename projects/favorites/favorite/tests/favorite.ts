import * as anchor from "@coral-xyz/anchor";
import { Program } from "@coral-xyz/anchor";
import { Favorite } from "../target/types/favorite";
import { assert } from 'chai';

describe("favorite", () => {
  // Configure the client to use the local cluster.
  const provider = anchor.AnchorProvider.env();
  anchor.setProvider(provider);
  const program = anchor.workspace.favorite as Program<Favorite>;
  const number = new anchor.BN(22)
  const color = "red"
  const hobbies = ['skiing', 'ball', 'game']



  // it("Is Created!", async () => {
  //   // Derive PDA with the same seeds as on-chain
  //   const userPubkey = provider.wallet.publicKey;
  //   const [favoritesPda] = anchor.web3.PublicKey.findProgramAddressSync(
  //     [Buffer.from('favorites'), userPubkey.toBuffer()],
  //     program.programId);

  //   // Invoke initialize with required args and accounts
  //   const tx = await program.methods
  //     .initialize(number, color, hobbies)
  //     .accounts({
  //       user: userPubkey,
  //     })
  //     .rpc();
  //   console.log("Your transaction signature", tx);

  //   // Fetch data from PDA and assert
  //   const dataFromPad = await program.account.favorite.fetch(favoritesPda);
  //   // And make sure it matches!
  //   assert.equal(dataFromPad.color, color);
  //   // A little extra work to make sure the BNs are equal
  //   assert.equal(dataFromPad.number.toString(), number.toString());
  //   // And check the hobbies too
  //   assert.deepEqual(dataFromPad.hobbies, hobbies);
  // });

  it("is updated!",async () => {
    const userPubkey = provider.wallet.publicKey;
    const [favoritesPda] = anchor.web3.PublicKey.findProgramAddressSync(
      [Buffer.from('favorites'), userPubkey.toBuffer()],
      program.programId);

    const tx = await program.methods
      .update(number, color, hobbies)
      .accounts({
        user: userPubkey,
      })
      .rpc();
    console.log("Your transaction signature", tx);

    // Fetch data from PDA and assert
    const dataFromPad = await program.account.favorite.fetch(favoritesPda);
    // And make sure it matches!
    assert.equal(dataFromPad.color, color);
    // A little extra work to make sure the BNs are equal
    assert.equal(dataFromPad.number.toString(), number.toString());
    // And check the hobbies too
    assert.deepEqual(dataFromPad.hobbies, hobbies);
  })
});
