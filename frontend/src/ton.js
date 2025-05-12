import { waitField } from './effects'

let connect = null
let connected = false

function getConnect() {
  return new Promise((resolve, reject) => {
    if (connect)
      return resolve({ tonConnect: connect, tonConnected: connected })
    waitField(() => window.TON_CONNECT_UI)
      .then(() => {
        const conn = new TON_CONNECT_UI.TonConnectUI({
          manifestUrl: `${location.origin}/tonconnect-manifest.json`
        })
        conn.setConnectRequestParameters({
          value: { tonProof: TG_USER.toString() }
        })
        conn.connectionRestored.then((restored) => {
          connected = restored
          connect = conn
          resolve({ tonConnect: connect, tonConnected: connected })
        })
      })
      .catch(reject)
  })
}

function setConnected(on) {
  connected = on
}

export default { getConnect, setConnected }
