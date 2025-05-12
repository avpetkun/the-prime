<script>
  import API from './apiAdmin'

  export let loader

  let rewardUserID = ''
  let rewardPoints = ''

  function onRewardUser() {
    loader(API.rewardUserPoints(rewardUserID, rewardPoints)).then(() => {
      console.info(`user ${rewardUserID} rewarded`)
      Telegram.WebApp.showAlert(`user ${rewardUserID} rewarded`)

      rewardUserID = ''
      rewardPoints = ''
    })
  }
</script>

<div class="block pop-page">
  <input type="text" placeholder="tg-user-id" bind:value={rewardUserID} />
  <input type="text" placeholder="points" bind:value={rewardPoints} />
  <button on:click={onRewardUser}>Reward user</button>
</div>

<style>
  .block {
    background: var(--tg-theme-section-bg-color, #fff);
    border-radius: 12px;
    padding: 16px;
    width: 100%;
    display: flex;
    flex-direction: column;
  }

  input {
    padding: 4px 4px;
    font-size: 14px;
    margin-bottom: 8px;
  }

  button {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 500;
    font-size: 14px;
    border-radius: 6px;
    padding: 10px;
    color: #fff;
    background: linear-gradient(90deg, #3b8cff 0%, #867cff 100%);
  }
</style>
