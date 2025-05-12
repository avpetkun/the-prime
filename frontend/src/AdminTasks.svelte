<script>
  import API from './apiAdmin'

  import PrimeIcon from '@/prime_26.png'
  import Icon from './AdminIcon.svelte'
  import Task from './Task.svelte'

  export let loader
  export let chats = [] // { chatID: 123, title: '', link: '' }
  export let tasks = []

  let showHidden = false

  const taskTypes = [
    'free',
    'invite',
    'join',
    'free_link',
    'partner_event',
    'partner_check',
    'ton_connect',
    'ton_disconnect',
    'ton_deposit',
    'stars_deposit',
    'ads_gram_task',
    'ads_gram_rewarded',
    'tapp_ads',
    'monetag-link',
    'monetag-banner'
  ]

  function getWebhook(task) {
    if (task.type == 'tapp_ads') {
      return `https://webhook.getprime.me/api/v1/tappads/${task.actionPartnerHook}`
    }
    return `https://webhook.getprime.me/api/v1/reward/${task.actionPartnerHook}?userid=123456789`
  }
  function webhookCopy(task) {
    navigator.clipboard.writeText(getWebhook(task))
  }

  function canIntervalType(t) {
    return !['join', 'free_link', 'ton_connect', 'ton_disconnect'].includes(t)
  }

  function canPartnerWebhook(t) {
    return [
      'partner_event',
      'ads_gram_task',
      'ads_gram_rewarded',
      'tapp_ads',
      'monetag-link',
      'monetag-banner'
    ].includes(t)
  }

  function newTask() {
    return {
      id: 0,
      type: '',
      name: '',
      desc: '',
      icon: '',
      hidden: true,
      premium: false,
      interval: 0,
      pending: 0,
      points: 100,
      weight: 0,
      maxClicks: 0,
      actionLink: '',
      actionChatID: 0,
      actionTonAmount: 0.5,
      actionStarsAmount: 1,
      actionStarsTitle: '',
      actionStarsDesc: '',
      actionStarsItem: '',
      actionPartnerHook: '',
      actionPartnerMatch: '',
      actionAdsGramBlockID: '',
      actionTappAdsToken: ''
    }
  }
  function copyTask(t) {
    return JSON.parse(JSON.stringify(t))
  }

  let edit = null

  function cancelNewTask() {
    tasks = tasks.slice(1)
    edit = null
  }

  function selectTask(t) {
    // prevent select duplicate
    if (edit && t.id == edit.id) return
    // if selected new task then delete new task
    if (edit && !edit.id) cancelNewTask()
    edit = copyTask(t)
  }

  function createTask() {
    // cancel new task if already creating
    if (edit && !edit.id) return cancelNewTask()
    edit = newTask()
    tasks = [edit].concat(tasks)
  }

  function saveTask() {
    loader(API.taskSave(edit)).then((formatted) => {
      if (formatted.id) {
        const i = tasks.findIndex((t) => t.id == edit.id)
        tasks[i] = formatted
        edit = null
      }
      tasks.sort((a, b) => {
        if (a.premium != b.premium) return b.premium ? 1 : -1
        return a.weight < b.weight ? 1 : -1
      })
      tasks = tasks
    })
  }
  function delTask(task) {
    if (!task.id) return cancelNewTask()
    function exec() {
      loader(API.taskDelete(task.id)).then(() => {
        const i = tasks.findIndex((t) => t.id == task.id)
        tasks.splice(i, 1)
        tasks = tasks
      })
    }
    if (window.DEBUG) return confirm('Delete task?') && exec()
    Telegram.WebApp.showConfirm('Delete task?', (ok) => ok && exec())
  }

  function openChat(chatID) {
    const chat = chats.find((c) => c.chatID == chatID)
    if (chat) window.open(chat.link)
  }

  const secretSize = 64
  function generateSecret() {
    const rand = new Uint8Array(secretSize)
    window.crypto.getRandomValues(rand)
    const alp = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    return Array.from(rand, (byte) => alp[byte % secretSize]).join('')
  }

  function enLocale(textWithLocals) {
    for (let line of textWithLocals.split('\n')) {
      if (line.startsWith('en:')) return line.slice(3).trim()
    }
    return textWithLocals
  }
</script>

<button class="add-task material dark pop-page" on:click={createTask}>
  {edit && !edit.id ? 'Close editor' : 'Create new task'}
</button>

{#if edit == null}
  <div class="hidden pop-page">
    <input
      id="show-hidden"
      type="checkbox"
      checked={showHidden}
      on:click={() => (showHidden = !showHidden)}
    />
    <label for="show-hidden">show hidden tasks</label>
  </div>
{/if}

<div class="tasks pop-page">
  <div></div>

  {#each tasks as t, id (id)}
    {#if edit && edit.id == t.id}
      <div class="editor">
        <div class="row">
          <span>Hide task</span>
          <input
            type="checkbox"
            checked={edit.hidden}
            on:click={(edit.hidden = !edit.hidden)}
          />
        </div>
        <div>
          <span>Max clicks count (now {edit.nowClicks})</span>
          <input type="number" bind:value={edit.maxClicks} />
        </div>
        <div>
          <span>Icon</span>
          <Icon bind:value={edit.icon} placeholder={PrimeIcon} />
        </div>
        <div>
          <span>Title</span>
          <textarea bind:value={edit.name} />
        </div>
        <div>
          <span>Description</span>
          <textarea bind:value={edit.desc} />
        </div>
        <div class="row">
          <span>Premium task</span>
          <input
            type="checkbox"
            checked={edit.premium}
            on:click={(edit.premium = !edit.premium)}
          />
        </div>
        <div>
          <span>Points</span>
          <input type="number" bind:value={edit.points} />
        </div>
        <div>
          <span>Weight</span>
          <input type="number" bind:value={edit.weight} />
        </div>
        <div>
          <span>Max pending time (in seconds)</span>
          <input type="number" bind:value={edit.pending} />
        </div>
        <div>
          <span>Type</span>
          <select bind:value={edit.type}>
            <option value="" hidden>Select task type</option>
            {#each taskTypes as tt}
              <option value={tt}>{tt}</option>
            {/each}
          </select>
        </div>
        {#if canIntervalType(edit.type)}
          <div>
            <span>Repeat interval (in seconds)</span>
            <input type="number" bind:value={edit.interval} />
          </div>
        {/if}
        {#if edit.type == 'free_link'}
          <div>
            <span>Link</span>
            <input type="text" bind:value={edit.actionLink} />
          </div>
        {/if}
        {#if edit.type == 'join'}
          <div>
            <span>Join Chat</span>
            <div class="input-button">
              <select bind:value={edit.actionChatID}>
                <option value={0} hidden>Select chat</option>
                {#each chats as c}
                  <option value={c.chatID}>{c.title}</option>
                {/each}
              </select>
              <button
                class="material dark"
                on:click={() => openChat(edit.actionChatID)}
              >
                open
              </button>
            </div>
          </div>
          <div>
            <span>Chat link</span>
            <input type="text" bind:value={edit.actionLink} />
          </div>
        {/if}
        {#if edit.type == 'ton_deposit'}
          <div>
            <span>Ton amount (with point)</span>
            <input type="number" bind:value={edit.actionTonAmount} />
          </div>
        {/if}
        {#if edit.type == 'stars_deposit'}
          <div>
            <span>Stars amount</span>
            <input type="number" bind:value={edit.actionStarsAmount} />
          </div>
          <div>
            <span>Stars invoice title</span>
            <input type="text" bind:value={edit.actionStarsTitle} />
          </div>
          <div>
            <span>Stars invoice description</span>
            <input type="text" bind:value={edit.actionStarsDesc} />
          </div>
          <div>
            <span>Stars invoice price item</span>
            <input type="text" bind:value={edit.actionStarsItem} />
          </div>
        {/if}
        {#if edit.type == 'partner_event' || edit.type == 'partner_check'}
          <div>
            <span>Partner app link</span>
            <textarea bind:value={edit.actionLink} />
          </div>
        {/if}
        {#if edit.type == 'partner_check'}
          <div>
            <span>Partner check link</span>
            <textarea bind:value={edit.actionPartnerHook} />
          </div>
          <div>
            <span>Partner expected response</span>
            <textarea bind:value={edit.actionPartnerMatch} />
          </div>
        {/if}
        {#if canPartnerWebhook(edit.type)}
          <div>
            <span>Partner API token</span>
            <div class="input-button">
              <input type="text" bind:value={edit.actionPartnerHook} />
              <button
                class="material dark"
                on:click={() => (edit.actionPartnerHook = generateSecret())}
              >
                generate
              </button>
            </div>
          </div>
          <div>
            <span>Webhook link</span>
            <div class="input-button">
              <textarea type="text" readonly value={getWebhook(edit)} />
              <button class="material dark" on:click={() => webhookCopy(edit)}>
                copy
              </button>
            </div>
          </div>
        {/if}
        {#if edit.type == 'tapp_ads'}
          <div>
            <span>TappAds auth token</span>
            <input type="text" bind:value={edit.actionTappAdsToken} />
          </div>
        {/if}
        {#if edit.type == 'ads_gram_task' || edit.type == 'ads_gram_rewarded'}
          <div>
            <span>AdsGram block-id</span>
            <input type="text" bind:value={edit.actionAdsGramBlockID} />
          </div>
        {/if}

        <div class="buttons">
          <button class="material dark" on:click={saveTask}>Save task</button>
          <button class="material dark delete" on:click={() => delTask(edit)}>
            Delete task
          </button>
        </div>
      </div>
    {:else if !t.hidden || t.hidden == showHidden}
      <Task
        premium={t.premium}
        icon={t.icon}
        name={enLocale(t.name)}
        desc={enLocale(t.desc)}
        points={t.points}
        on:start={() => selectTask(t)}
      />
    {/if}
  {/each}
</div>

<style>
  * {
    color: var(--tg-theme-text-color, #000);
  }

  .hidden {
    margin-bottom: 16px;
    font-size: 14px;
  }

  .add-task {
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

  button {
    background: #cae3f3;
    color: #168dcd;
    font-size: 14px;
    font-weight: 500;
    border-radius: 10px;
    overflow: hidden;
  }
  button.delete {
    background: #f3caca;
    color: #cd1616;
  }

  .tasks {
    display: flex;
    flex-direction: column;
  }

  .editor {
    display: flex;
    flex-direction: column;
    background: #fff;
    padding: 16px;
    border-top: 1px solid #999999;
    border-bottom: 1px solid #999999;
  }

  .editor:last-child {
    border-bottom: none;
    border-radius: 0 0 12px 12px;
  }
  .editor:nth-child(2) {
    border-top: none;
    border-radius: 12px 12px 0 0;
  }
  .editor:last-child:nth-child(2) {
    border-top: none;
    border-bottom: none;
    border-radius: 12px;
  }

  .editor,
  .editor input,
  .editor textarea,
  .editor select {
    width: 100%;
  }
  .editor > div {
    display: flex;
  }
  .editor > div:not(:last-child) {
    margin-bottom: 6px;
  }
  .editor > div:not(.row) {
    flex-direction: column;
  }
  .editor span {
    font-size: 14px;
    font-weight: 500;
    margin-bottom: 2px;
  }
  .editor input[type='checkbox'] {
    width: auto;
    margin-left: 10px;
  }

  .editor .input-button {
    display: flex;
    flex-direction: row;
  }
  .editor .input-button button {
    flex-shrink: 0;
    margin-left: 5px;
  }

  .editor > div.buttons {
    margin-top: 20px;
    display: flex;
    flex-direction: row;
  }
  .editor .buttons button {
    width: 50%;
    height: 30px;
  }
</style>
