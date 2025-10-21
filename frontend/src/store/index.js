import { createStore } from 'vuex'
import auth from './modules/auth'
import animals from './modules/animals'
import adoptions from './modules/adoptions'
import schedules from './modules/schedules'
import documents from './modules/documents'
import finances from './modules/finances'
import users from './modules/users'

export default createStore({
  modules: {
    auth,
    animals,
    adoptions,
    schedules,
    documents,
    finances,
    users
  }
})
