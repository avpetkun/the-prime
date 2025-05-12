function exec(json, method, url, body) {
  return new Promise((resolve, reject) => {
    fetch(url, {
      method: method || 'GET',
      headers: {
        Authorization: AUTH,
        'Content-Type': 'application/json'
      },
      body: body ? JSON.stringify(body) : undefined
    })
      .then((res) => {
        if (res.status >= 300)
          res
            .json()
            .then((msg) => reject(msg.message || msg))
            .catch(() => reject(`error status ${res.status}`))
        else if (!json) resolve()
        else res.json().then(resolve)
      })
      .catch(reject)
  })
}

function getOverview() {
  return exec(true, 'GET', '/api/v1/overview-admin')
}

function productSave(product) {
  return exec(true, 'POST', '/api/v1/products', product)
}

function productDelete(productID) {
  return exec(false, 'DELETE', `/api/v1/products/${productID}`)
}

function taskSave(task) {
  return exec(true, 'POST', `/api/v1/tasks`, task)
}

function taskDelete(taskID) {
  return exec(false, 'DELETE', `/api/v1/tasks/${taskID}`)
}

function rewardUserPoints(userID, points) {
  return exec(
    false,
    'POST',
    `/api/v1/reward-user?user=${userID}&points=${points}`
  )
}

export default {
  getOverview,
  productSave,
  productDelete,
  taskSave,
  taskDelete,
  rewardUserPoints
}
