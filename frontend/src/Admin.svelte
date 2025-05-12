<script>
  import API from './apiAdmin'

  import AdminTasks from './AdminTasks.svelte'
  import AdminShop from './AdminShop.svelte'
  import AdminManualReward from './AdminManualReward.svelte'

  import TgStarIcon from '@/tgstar.svg'

  let pages = [
    { name: 'Tasks', component: AdminTasks },
    { name: 'Shop', component: AdminShop },
    { name: 'Manual reward', component: AdminManualReward }
  ]
  let page = null

  window.Reload = () => (page = null)

  let products = []
  let chats = []
  let tasks = []

  let loading = false
  function loader(promise) {
    loading = true
    return promise
      .catch((error) => {
        console.error(error)
        Telegram.WebApp.showAlert(error)
      })
      .finally(() => (loading = false))
  }

  loader(API.getOverview()).then((o) => {
    products = o.products
    chats = o.chats
    tasks = o.tasks
  })
</script>

{#if loading}
  <div class="loader">
    <img src={TgStarIcon} alt="" />
  </div>
{/if}

<div class="admin" class:loading>
  {#if page && page.component}
    <svelte:component
      this={page.component}
      {loader}
      {products}
      {tasks}
      {chats}
    />
  {:else}
    <div class="pages pop-page">
      {#each pages as p}
        <button class="material dark" on:click={() => (page = p)}>
          {p.name}
        </button>
      {/each}
    </div>
  {/if}
</div>

<style>
  .admin {
    position: relative;
    width: 100%;
    min-height: 100%;
    padding: 0 16px;
    transition: filter 0.3s linear;
  }
  .admin.loading {
    filter: blur(4px);
  }

  .loader {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    z-index: 10;
  }
  .loader img {
    position: absolute;
    top: 50%;
    left: calc(50% - 20px);
    width: 40px;
    height: 40px;
    animation: bounce-star 0.3s infinite ease-in-out;
    z-index: 10;
  }

  .pages {
    width: 100%;
    height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }
  .pages button {
    margin-bottom: 8px;
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
