import Generic from './generic-store.js'
import PlayerStore from './player-store.js'
import HistoryStore from './history-store.js'

export default new Vuex.Store({
  modules: {
    gen: Generic,
    ps: PlayerStore,
    hs: HistoryStore,
  }
})
