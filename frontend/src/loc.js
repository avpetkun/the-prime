const en = {
  claim: 'Claim',
  noTasks: 'No tasks',
  points: 'points',
  top: {
    title: 'Get Telegram Premium!',
    desc1: 'Complete <b>tasks</b> to earn Prime Points!',
    desc2: 'Use <b>points</b> in the shop to get <b>digital goods</b>'
  },
  blocks: {
    premium: 'PREMIUM TASKS',
    claim: 'CLAIM TASKS',
    simple: 'SIMPLE TASKS',
    pending: 'PENDING TASKS',
    done: 'DONE TASKS'
  },
  tabs: {
    active: 'Active',
    pending: 'Pending',
    done: 'Done'
  }
}

const ru = {
  claim: 'Забрать',
  noTasks: 'Нет задач',
  points: 'баллов',
  top: {
    title: 'Получи Telegram Premium!',
    desc1: 'Зарабатывай баллы выполняя <b>задачи</b>',
    desc2: 'Обменивай баллы на <b>награды</b> в магазине!'
  },
  blocks: {
    premium: 'ПРЕМИУМ ЗАДАЧИ',
    claim: 'ПОЛУЧИ ЗАДАЧИ',
    simple: 'ОБЫЧНЫЕ ЗАДАЧИ',
    pending: 'ЗАДАЧИ НА ПРОВЕРКЕ',
    done: 'ВЫПОЛНЕННЫЕ ЗАДАЧИ'
  },
  tabs: {
    active: 'Активные',
    pending: 'Проверка',
    done: 'Сделано'
  }
}

const locs = { en, ru }

export default () => locs[window.LANG] || locs.en
