import Vue from 'vue';
import Router from 'vue-router';
import Index from './pages/Index.vue';
import MainNavbar from './layout/MainNavbar.vue';
import MainFooter from './layout/MainFooter.vue';
import Content from './pages/Content.vue';

Vue.use(Router);

export default new Router({
    linkExactActiveClass: 'active',
    routes: [
        {
            path: '/',
            name: 'index',
            components: {default: Index, header: MainNavbar, footer: MainFooter},
            props: {
                header: {colorOnScroll: 400},
                footer: {backgroundColor: 'black'}
            }
        },
        {
            path: '/:id',
            name: "content",
            components: {default: Content, header: MainNavbar, footer: MainFooter},
            props: {
                default: true,
                header: {colorOnScroll: 400},
                footer: {backgroundColor: 'black'}
            }
        }
    ],
    scrollBehavior: to => {
        if (to.hash) {
            return {selector: to.hash};
        } else {
            return {x: 0, y: 0};
        }
    }
});
