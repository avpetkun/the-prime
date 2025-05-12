function exec(json, method, url, body) {
  return new Promise((resolve, reject) => {
    fetch(url, {
      method: method || 'GET',
      headers: {
        Authorization: AUTH
      },
      body: body ? JSON.stringify(body) : undefined
    })
      .then((res) => {
        if (res.status >= 300) reject()
        else if (!json) resolve()
        else res.json().then(resolve)
      })
      .catch(reject)
  })
}

function getOverview(start) {
  return exec(true, 'GET', `/api/v1/overview?start=${start}`)
}

function getTasksEvents(startTime) {
  // [
  //   { id: 26, subID: 0, status: 'active' },
  //   { id: 27, subID: 5, status: 'claim' }
  // ]
  return exec(true, 'GET', `/api/v1/tasks/events?from=${startTime}`)
}

function taskStart(taskID, subID) {
  return exec(false, 'POST', `/api/v1/tasks/${taskID}/${subID}/start`)
}

function taskClaim(taskID, subID) {
  return exec(false, 'POST', `/api/v1/tasks/${taskID}/${subID}/claim`)
}

function getTaskStarsInvoice(taskID) {
  // '<invoice-link>'
  return exec(true, 'POST', `/api/v1/tasks/${taskID}/invoice/stars`)
}

function getTaskTonInvoice(taskID) {
  // {
  //   validUntil: 123,
  //   messages: [
  //     {
  //       address: '',
  //       amount: '',
  //       payload: ''
  //     }
  //   ]
  // }
  return exec(true, 'POST', `/api/v1/tasks/${taskID}/invoice/ton`)
}

function getInviteMessage() {
  // "<invite-msg-id>"
  return exec(true, 'POST', `/api/v1/invite-message`)
}

function sendInit() {
  return exec(false, 'POST', `/api/v1/init`)
}

function productClaim(productID) {
  return exec(true, 'POST', `/api/v1/products/${productID}/claim`)
}

export default {
  getOverview,
  getTasksEvents,
  taskStart,
  taskClaim,
  getTaskStarsInvoice,
  getTaskTonInvoice,
  getInviteMessage,
  sendInit,
  productClaim
}
