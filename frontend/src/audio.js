const audioContext = new (window.AudioContext || window.webkitAudioContext)()

function loadSound(filename) {
  const sound = { buf: null }
  const ajax = new XMLHttpRequest()
  ajax.open('GET', filename, true)
  ajax.responseType = 'arraybuffer'
  ajax.onload = () =>
    audioContext.decodeAudioData(ajax.response, (buf) => (sound.buf = buf))
  ajax.send()
  return sound
}

function playSound(sound, volume) {
  if (sound.buf) {
    const source = audioContext.createBufferSource()
    if (source) {
      source.buffer = sound.buf

      if (!source.start) source.start = source.noteOn
      if (source.start) {
        const gain = audioContext.createGain()
        gain.gain.value = volume
        source.connect(gain)
        gain.connect(audioContext.destination)
        source.start(0)
      }
    }
  }
}

let needInit = true

export default function (filename) {
  const sound = loadSound(filename)
  if (needInit) {
    needInit = false
    const initAudio = () => {
      playSound(sound, 0)
      window.removeEventListener('touchstart', initAudio, false)
      window.removeEventListener('click', initAudio, false)
    }
    window.addEventListener('touchstart', initAudio, false)
    window.addEventListener('click', initAudio, false)
  }
  return () => playSound(sound, 1)
}
