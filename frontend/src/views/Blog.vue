<template>
  <div>
    <!-- Blog Header -->
    <section class="blog-header">
      <v-container class="py-16">
        <v-row justify="center">
          <v-col cols="12" md="8" class="text-center">
            <h1 class="text-h3 font-weight-bold mb-4">
              健康數據研究與發展
            </h1>
            <p class="text-h6">
              探索健康大數據的最新發展與應用
            </p>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- Blog Posts -->
    <section class="blog-posts">
      <v-container class="py-16">
        <v-row>
          <v-col
            v-for="post in blogPosts"
            :key="post.title"
            cols="12"
            md="4"
            class="mb-6"
          >
            <v-card
              class="h-100"
              elevation="2"
              :to="post.link"
            >
              <v-img
                :src="post.image"
                height="200"
                cover
                class="align-end"
              >
                <v-card-title class="text-white text-h5">
                  {{ post.title }}
                </v-card-title>
              </v-img>

              <v-card-text>
                <div class="text-subtitle-1 mb-2">
                  {{ post.date }}
                </div>
                <p class="text-body-1">
                  {{ post.excerpt }}
                </p>
              </v-card-text>

              <v-card-actions>
                <v-btn
                  color="primary"
                  variant="text"
                  :to="post.link"
                >
                  閱讀更多
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const blogPosts = ref([
  {
    title: '醫療科技未來發展',
    date: '2024-03-15',
    excerpt: '探索人工智慧在醫療領域的應用，以及如何改變未來的醫療服務模式。',
    image: 'https://picsum.photos/800/400?random=1',
    link: '/blog/ai-healthcare'
  },
  {
    title: '預防保健小貼士',
    date: '2024-03-10',
    excerpt: '了解如何通過日常習慣和定期檢查來維護健康，預防疾病。',
    image: 'https://picsum.photos/800/400?random=2',
    link: '/blog/preventive-care'
  },
  {
    title: '心理健康意識',
    date: '2024-03-05',
    excerpt: '探討心理健康的重要性，以及如何在現代生活中保持心理平衡。',
    image: 'https://picsum.photos/800/400?random=3',
    link: '/blog/mental-health'
  }
])
</script>

<style scoped>
.blog-header {
  background: linear-gradient(135deg, var(--v-primary-base) 0%, var(--v-primary-darken2) 100%);
  color: white;
}

.blog-posts {
  background-color: #f8f9fa;
}

.v-card {
  transition: transform 0.3s ease;
}

.v-card:hover {
  transform: translateY(-5px);
}
</style> 