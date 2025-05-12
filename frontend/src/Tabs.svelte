<script>
  import './tabs.scss'
  import { startFollow } from './effects'

  import TabActiveIcon from '@/tab_active.svg'
  import TabPendingIcon from '@/tab_pending.svg'
  import TabDoneIcon from '@/tab_done.svg'

  import getLoc from './loc'
  const loc = getLoc()

  export let value = 0

  export let countActive = 0
  export let countPending = 0
  export let countDone = 0

  const tabs = [
    { icon: TabActiveIcon, name: loc.tabs.active },
    { icon: TabPendingIcon, name: loc.tabs.pending },
    { icon: TabDoneIcon, name: loc.tabs.done }
  ]

  let activeAnim = 0
  let pendingAnim = 0
  let doneAnim = 0

  function animateValues() {
    const active0 = activeAnim
    const pending0 = pendingAnim
    const done0 = doneAnim
    startFollow(500, (follow) => {
      activeAnim = follow(active0, countActive)
      pendingAnim = follow(pending0, countPending)
      doneAnim = follow(done0, countDone)
    })
  }
  $: animateValues(countActive, countPending, countDone)
</script>

<div class="tabs">
  {#each tabs as t, i}
    <button
      class="material"
      class:dark={value == i}
      on:click={() => (value = i)}
    >
      <svg>
        <use xlink:href={t.icon + '#i'}></use>
      </svg>
      <span>{t.name}</span>
    </button>
  {/each}
  <div
    class="badge active"
    class:show={activeAnim}
    class:zoom={countActive > activeAnim}
  >
    +{activeAnim}
  </div>
  <div
    class="badge pending"
    class:show={pendingAnim}
    class:zoom={countPending > pendingAnim}
  >
    +{pendingAnim}
  </div>
  <div
    class="badge done"
    class:show={doneAnim}
    class:zoom={countDone > doneAnim}
  >
    +{doneAnim}
  </div>
</div>
