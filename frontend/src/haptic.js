const haptic = Telegram.WebApp.HapticFeedback

export default {
  lightInpact: () => haptic.impactOccurred('light'),
  notifySuccess: () => haptic.notificationOccurred('success'),
  notifyError: () => haptic.notificationOccurred('error'),
  notifyWarn: () => haptic.notificationOccurred('warning')
}
