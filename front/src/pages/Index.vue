<template>
  <div>
    <div class="page-header clear-filter" filter-color="orange">
      <div class="container">
        <el-carousel
          :height="carouselHeight"
          :interval="5000"
          arrow="hover"
          indicator-position="none"
          @change="changeHeight()"
        >
          <el-carousel-item v-for="(image,index) in homeContext.carousel" :key="index">
            <img
              id="banner-img"
              :src="'https://sniweb.shouting.feedia.co/php'+image.image"
              :alt="image.id"
              width="100%"
              @click="openLink(image.url)"
              style="cursor:pointer;"
            />
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
          <div class>
            <div v-html="homeContext.page_data.html_context"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { Parallax } from "@/components";
import { mapActions, mapGetters } from "vuex";
import { Carousel, CarouselItem } from "element-ui";
import MarqueeText from "vue-marquee-text-component";

export default {
  name: "index",
  bodyClass: "index-page",
  components: {
    Parallax,
    "el-carousel": Carousel,
    "el-carousel-item": CarouselItem,
    MarqueeText
  },
  computed: {
    ...mapGetters(["homeContext", "headerNav", "homeContent"]),
    marquee() {
      let marqueeList = [];

      if (this.homeContext.marquee) {
        this.homeContext.marquee.forEach(value => {
          marqueeList.push(
            `<span class="ml-5" style="color:${value.color}">${value.text}</span>`
          );
        });
      }

      return marqueeList.join(
        "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"
      );
    }
  },
  data() {
    return {
      carouselHeight: "300px"
    };
  },
  watch: {
      homeContext(){
          this.changeHeight();
      }
  },
  mounted() {
    const that = this;
    window.onresize = () => {
      return (() => {
        that.changeHeight();
      })();
    };
  },
  created() {},
  methods: {
    ...mapActions([
      "actionUpdateContext",
      "actionHeaderNav",
      "actionSetContent"
    ]),
    openLink(url) {
      window.location.href = url;
    },
    changeHeight() {
      if (document.getElementById("banner-img")) {
        this.carouselHeight =
          document.getElementById("banner-img").clientHeight + "px";
      }
    }
  }
};
</script>
<style scoped>
.page-header {
}
</style>
