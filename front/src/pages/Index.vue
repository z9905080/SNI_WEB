<template>
    <div>
        <div class="page-header clear-filter" filter-color="orange">

            <div class="container">
                <el-carousel :interval="5000" arrow="hover" indicator-position="none">
                    <el-carousel-item v-for="(image,index) in homeContext.carousel" :key="index">
                           <img :src="'https://sniweb.shouting.feedia.co/php'+image.image" :alt="image.id" height="100%" @click="openLink(image.url)" style="cursor:pointer;"/>
                    </el-carousel-item>
                </el-carousel>
            </div>
        </div>
        <div class="main">
            <div class="section section-images">
                <div class="container">
                    <marquee-text :repeat="2">
                        <div v-html="marquee"></div>
                    </marquee-text>
                    <div class="">
                        <div v-html="homeContext.page_data.html_context"></div>
                    </div>
                </div>
            </div>
        </div>

    </div>
</template>
<script>
    import {Parallax} from '@/components';
    import {mapActions, mapGetters} from "vuex";
    import {Carousel, CarouselItem} from "element-ui"
    import MarqueeText from 'vue-marquee-text-component'


    export default {
        name: 'index',
        bodyClass: 'index-page',
        components: {
            Parallax,
            "el-carousel": Carousel,
            "el-carousel-item": CarouselItem,
            MarqueeText
        },
        computed: {
            ...mapGetters([
                'homeContext',
                'headerNav',
                'homeContent'
            ]),
            marquee(){
                let marqueeList = [];

                if (this.homeContext.marquee) {
                    this.homeContext.marquee.forEach(value => {
                        marqueeList.push(`<span class="ml-5" style="color:${value.color}">${value.text}</span>`)
                    });
                }

                return marqueeList.join("&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;");
            }
        },
        created() {

        },
        methods: {
            ...mapActions([
                'actionUpdateContext',
                'actionHeaderNav',
                'actionSetContent'
            ]),
            openLink(url){
                window.location.href = url;
            }
        }
    };
</script>
<style scoped>
    .page-header{

    }

</style>
