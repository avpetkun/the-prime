<script>
  import './shop.scss'

  import API from './apiAdmin'

  import TgStarIcon from '@/tgstar.svg'
  import TgPrimeIcon from '@/prime_96.png'

  export let loader
  export let products = []

  function typeIcon(type) {
    return {
      tg_stars: TgStarIcon,
      tg_premium: TgPrimeIcon
    }[type]
  }

  function fmtNumber(n) {
    return n.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
  }

  function remaining(p) {
    let r = p.price
    if (r > 9999) r = `${Math.floor(r / 1000)}k`
    return `Earn ${r} points to claim`
  }

  function select(p) {
    product = JSON.parse(JSON.stringify(p))
  }

  const productTypes = ['tg_stars', 'tg_premium']

  let product = null

  function addProduct() {
    if (product) return (product = null)
    product = {
      id: 0,
      type: 'tg_stars',
      name: '',
      price: 0,
      amount: 0,
      badge: ''
    }
  }

  function saveProduct() {
    loader(API.productSave(product)).then((p) => {
      const i = products.findIndex((v) => v.id == p.id)
      if (i == -1) {
        products.push(p)
      } else {
        products[i] = p
      }
      products = products
      product = null
    })
  }

  function delProduct() {
    function exec() {
      loader(API.productDelete(product.id)).then(() => {
        const i = products.findIndex((v) => v.id == product.id)
        if (i != -1) {
          products.splice(i, 1)
          products = products
        }
        product = null
      })
    }
    if (window.DEBUG) return confirm('Delete product?') && exec()
    Telegram.WebApp.showConfirm('Delete product?', (ok) => ok && exec())
  }
</script>

<div class="shop pop-page">
  <button class="add-product material dark" on:click={addProduct}>
    {product ? 'Close editor' : 'Add new product'}
  </button>

  {#if product}
    <div class="edit-buttons">
      <button class="material dark" on:click={saveProduct}>Save</button>
      {#if product.id}
        <button class="material dark delete" on:click={delProduct}>
          Delete
        </button>
      {/if}
    </div>

    <div class="edit-fields">
      <div>
        <span>Product type</span>
        <select bind:value={product.type}>
          <option value={0} hidden>Select type</option>
          {#each productTypes as t}
            <option value={t}>{t}</option>
          {/each}
        </select>
      </div>
      <div>
        <span>Product name</span>
        <input bind:value={product.name} />
      </div>
      <div>
        <span>Price</span>
        <input
          type="number"
          min="0"
          step="1"
          pattern="\d+"
          bind:value={product.price}
        />
      </div>
      <div>
        <span>Amount</span>
        <input
          type="number"
          min="0"
          step="1"
          pattern="\d+"
          bind:value={product.amount}
        />
      </div>
      <div>
        <span>Badge</span>
        <input bind:value={product.badge} />
      </div>
    </div>
  {:else}
    <div class="products">
      {#each products as p}
        <button class="product material" on:click={() => select(p)}>
          <div class="icon">
            <img src={typeIcon(p.type)} alt="" />
          </div>
          <div class="content">
            <div class="name">
              <span>{p.name}</span>
              {#if p.badge}
                <span class="badge">{p.badge}</span>
              {/if}
            </div>
            <span class="info">{remaining(p)}</span>
          </div>
          {#if p.price}
            <div class="points active">{fmtNumber(p.price)}</div>
          {/if}
        </button>
      {/each}
    </div>
  {/if}
</div>

<style>
  .shop .products {
    padding: 0;
  }

  input,
  select {
    padding: 4px 4px;
    font-size: 14px;
  }

  .add-product {
    margin-bottom: 16px;
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

  .edit-buttons {
    margin-bottom: 16px;
    display: flex;
  }

  .edit-buttons button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    flex-grow: 1;
    font-weight: 500;
    font-size: 14px;
    border-radius: 6px;
    padding: 10px;
    color: #fff;
    background: rgb(112, 195, 120);
  }
  .edit-buttons button.delete {
    color: #fff;
    background: #e27575;
  }
  .edit-buttons button:nth-child(2) {
    margin-left: 8px;
  }

  .edit-fields {
    background: var(--tg-theme-section-bg-color, #fff);
    border-radius: 12px;
    padding: 16px;

    display: flex;
    flex-direction: column;
  }
  .edit-fields > div {
    display: flex;
    flex-direction: column;
    color: var(--tg-theme-text-color, #000);
    font-weight: 500;
    font-size: 16px;
  }
  .edit-fields > div > span:first-child {
    margin: 4px 0 2px 0;
  }
</style>
