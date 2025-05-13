import { createStore } from 'vuex';
import auth from './auth';
import animals from './animals';
import adoptions from './adoptions';
import schedule from './schedule';
import documents from './documents';
import finances from './finances';

export default createStore({
  modules: {
    animals,
    adoptions,
    schedule,
    documents,
    finances,
    users: {},
    auth
  }
}); 