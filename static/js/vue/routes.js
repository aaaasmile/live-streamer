import Dashboard from './views/dashboard.js'
import OSView from './views/os.js'
import HistoryView from './views/history.js'

export default [
  { path: '/', icon: 'dashboard', title: 'Dashboard', component: Dashboard },
  { path: '/os', icon: 'dashboard', title: 'OS', component: OSView },
  { path: '/history', icon: 'dashboard', title: 'History', component: HistoryView },
]