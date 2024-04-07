const routes = [
    { path: '/', component: httpVueLoader( 'components/configlist.vue' ) },
    { path: '/add', component: httpVueLoader( 'components/configadd.vue' ) },
];
const router = new VueRouter({
    routes
});
export default router