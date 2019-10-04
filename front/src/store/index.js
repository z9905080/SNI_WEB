const state = {
    homeContext: {
        web_title:'',
        page_data:{
            html_context:''
        }
    },
    headerNav: [],
    homeContent: {}
};

// getters
const getters = {
    homeContext: state => state.homeContext,
    headerNav: state => state.headerNav,
    homeContent: state => state.homeContent
};
// actions
const actions = {
    actionUpdateContext({ commit }, data) {
        commit('updateHomeContext', data)
    },
    actionHeaderNav({ commit }, data) {
        commit('updateHeaderNav', data);
    },
    actionSetContent({ commit }, data) {
        commit('updateContent', data);
    }
};

// mutations
const mutations = {
    updateHomeContext(state, payload) {
        state.homeContext = payload;
    },
    updateHeaderNav(state, payload) {
        state.headerNav = payload;
    },
    updateContent(state, payload) {
        state.homeContent = payload;
    }
};

export default {
    state,
    getters,
    actions,
    mutations,
};