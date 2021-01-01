export default {
  state: {
    errorText: '',
    msgText: '',
  },
  mutations: {
    errorText(state, msg) {
      state.errorText = msg
    },
    msgText(state, msg) {
      state.msgText = msg
    },
    clearErrorText(state) {
      if (state.errorText !== '') {
        state.errorText = ''
      }
    },
    clearMsgText(state) {
      if (state.msgText !== '') {
        state.msgText = ''
      }
    }
  }
}