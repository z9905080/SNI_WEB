<template>
    <navbar
            position="fixed"
            type="primary"
            :transparent="transparent"
            :color-on-scroll="colorOnScroll"
            menu-classes="ml-auto"
    >
        <template slot-scope="{ toggle, isToggled }">
            <router-link class="navbar-brand" to="/">
                {{homeContext.web_title}}
            </router-link>
        </template>
        <template slot="navbar-menu">
            <li class="nav-item">
                <a
                        class="nav-link"
                        href="/"
                >
                    <p>首頁</p>
                </a>
            </li>

            <drop-down
                    v-for="(nav,index) in headerNav.slice(1,5)" :key="index"
                    tag="li"
                    :title="nav.group_name"
                    class="nav-item"
            >
                <nav-link v-for="(child , i) in nav.page_content" :key="i" :to="'/'+ child.page_content_id">
                    {{child.page_name}}
                </nav-link>
            </drop-down>


            <drop-down
                    v-if="headerNav.length > 5"
                    tag="li"
                    title="其他"
                    class="nav-item"
            >
                <li class="dropdown-submenu"
                    v-for="(nav,index) in headerNav.slice(5)" :key="index"
                >
                    <div class="dropdown-item dropdown-toggle" @click="openSubMenu(nav.page_group_id)">
                        {{nav.group_name}}
                    </div>
                    <ul class="dropdown-menu" v-if="subMenuId === nav.page_group_id">
                        <nav-link v-for="(child , i) in nav.page_content" :key="i" :to="'/'+ child.page_content_id">
                            {{child.page_name}}
                        </nav-link>
                    </ul>
                </li>
            </drop-down>


        </template>
    </navbar>
</template>

<script>
    import {DropDown, NavbarToggleButton, Navbar, NavLink} from '@/components';
    import axios from 'axios'
    import {mapActions, mapGetters} from "vuex";

    export default {
        name: 'main-navbar',
        data() {
            return {
                subMenuId: undefined
            }
        },
        props: {
            transparent: Boolean,
            colorOnScroll: Number
        },
        components: {
            DropDown,
            Navbar,
            NavbarToggleButton,
            NavLink
        },
        computed: {
            ...mapGetters([
                'homeContext',
                'headerNav',
                'homeContent'
            ])
        },
        created() {
            axios.get("https://sniweb.shouting.feedia.co/php/GetNav.php", {
                'Cache-Control': 'max-age=3600'
            }).then(((res) => {
                this.actionHeaderNav(res.data.filter(nav => nav.page_content));
            }));

            axios.get("https://sniweb.shouting.feedia.co/php/GetHomePage.php", {
                'Cache-Control': 'max-age=3600'
            }).then(((res) => {
                this.actionUpdateContext(res.data);
            }));

        },
        methods: {
            ...mapActions([
                'actionUpdateContext',
                'actionHeaderNav',
                'actionSetContent'
            ]),
            openSubMenu(id) {
                this.subMenuId = id;
            }
        }
    };
</script>

<style scoped>

    .dropdown-submenu {
        position: relative;
    }

    .dropdown-submenu > .dropdown-menu {
        top: 0;
        left: -100%;
        margin-top: -6px;
    }

    .dropdown-submenu > .dropdown-menu:before {
        content: none;
    }
</style>
