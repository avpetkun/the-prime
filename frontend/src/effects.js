function waitField(f) {
  return new Promise((resolve) => {
    if (f()) return resolve(f())
    let i = setInterval(() => {
      if (f()) {
        clearInterval(i)
        resolve(f())
      }
    }, 100)
  })
}

function startFollow(ms, follow /* func((start,target,float) => val) */) {
  let startTime = null
  function step(timestamp) {
    if (!startTime) startTime = timestamp
    const progress = Math.min((timestamp - startTime) / ms, 1)
    follow((start, target, float) => {
      start += progress * (target - start)
      if (!float) start = Math.floor(start)
      return start
    })
    if (progress < 1) requestAnimationFrame(step)
  }
  requestAnimationFrame(step)
}

export { waitField, startFollow }
