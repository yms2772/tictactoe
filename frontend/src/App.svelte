<script lang="ts">
  import {
    CheckComputerWin,
    CheckUserWin,
    GetFields,
    SetComputerMark,
    SetUserMark
  } from "../wailsjs/go/main/App.js";

  let fields = [[{mark: '', willBeRemoved: false}]];

  async function updateFields() {
    fields = await GetFields();
  }

  async function setComputerMark(): Promise<void> {
    await SetComputerMark();
  }

  async function checkUserWin(): Promise<boolean> {
    return await CheckUserWin();
  }

  async function checkComputerWin(): Promise<boolean> {
    return await CheckComputerWin();
  }

  async function setUserMark(n: number): Promise<boolean> {
    return await SetUserMark(n);
  }

  async function click(n: number): Promise<void> {
    const marked = await setUserMark(n);
    if (!marked) return;
    await updateFields();
    const userWin = await checkUserWin();
    if (userWin) {
        await updateFields();
        return;
    }

    await setComputerMark();
    await updateFields();
    const computerWin = await checkComputerWin();
    if (computerWin) {
      await updateFields();
      return;
    }
  }

  updateFields();
</script>

<main>
  <div id="wrapper">
    {#each fields as row, rowIndex}
      {#each row as cell, colIndex}
        <button class="field {cell.mark === 'O' ? 'user' : 'computer'} {cell.willBeRemoved === true ? 'remove' : ''}" on:click={() => click(rowIndex * 3 + colIndex)}>
          {cell.mark}
        </button>
      {/each}
    {/each}
  </div>
</main>