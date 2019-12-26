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
              :src="'http://www.seicho-no-ie.org.tw/php'+image.image"
              :alt="image.id"
              width="100%"
              @load="changeHeight()"
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
          <div class v-if="homeContent">
            <div v-html="homeContent.html_context"></div>
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
import axios from "axios";
import MarqueeText from "vue-marquee-text-component";

export default {
  name: "index",
  bodyClass: "index-page",
  props: {
    id: String
  },
  components: {
    Parallax,
    "el-carousel": Carousel,
    "el-carousel-item": CarouselItem,
    MarqueeText
  },
  watch: {
    // 如果 `question` 发生改变，这个函数就会运行
    id: function(value, oldValue) {
      if (value !== oldValue) {
        this.homeContent.html_context = "";
        this.getContent();
      }
    },
    homeContext() {
      this.changeHeight();
    }
  },
  data() {
    return {
      carouselHeight: "300px"
    };
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
  mounted() {
    this.getContent();
    const that = this;
    window.onresize = () => {
      return (() => {
        that.changeHeight();
      })();
    };
  },
  methods: {
    ...mapActions([
      "actionUpdateContext",
      "actionHeaderNav",
      "actionSetContent"
    ]),
    getContent() {
      axios
        .get("http://www.seicho-no-ie.org.tw/php/GetContent.php", {
          params: {
            page_id: this.id
          }
        })
        .then(response => {
          this.actionSetContent(response.data);
        });
    },
    openLink(url) {
      window.location.href = url;
    },
    changeHeight() {
      if (document.getElementById("banner-img")) {
        if (document.getElementById("banner-img").clientHeight > 0) {
          this.carouselHeight =
            document.getElementById("banner-img").clientHeight + "px";
        }
      }
    }
  }
};
</script>
<style scoped>
.page-header {
}
</style>
