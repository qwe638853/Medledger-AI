<template>
  <div>
    <!-- Blog Header -->
    <section class="blog-header">
      <v-container class="py-16">
        <v-row justify="center">
          <v-col cols="12" md="8" class="text-center">
            <h1 class="text-h2 font-weight-bold mb-6">
              Our Blog
            </h1>
            <p class="text-h6">
              Stay updated with the latest healthcare insights and news
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
            :key="post.id"
            cols="12"
            md="4"
          >
            <v-card class="h-100">
              <v-img
                :src="post.image"
                height="200"
                cover
              ></v-img>
              <v-card-title>{{ post.title }}</v-card-title>
              <v-card-subtitle>{{ post.date }}</v-card-subtitle>
              <v-card-text>{{ post.excerpt }}</v-card-text>
              <v-card-actions>
                <v-btn
                  color="primary"
                  text
                >
                  Read More
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
    id: 1,
    title: 'The Future of Healthcare Technology',
    date: 'March 15, 2024',
    excerpt: 'Discover how emerging technologies are transforming the healthcare industry...',
    image: 'https://picsum.photos/seed/health1/800/400'
  },
  {
    id: 2,
    title: 'Preventive Healthcare Tips',
    date: 'March 10, 2024',
    excerpt: 'Learn about the best practices for maintaining good health and preventing illness...',
    image: 'https://picsum.photos/seed/health2/800/400'
  },
  {
    id: 3,
    title: 'Mental Health Awareness',
    date: 'March 5, 2024',
    excerpt: 'Understanding the importance of mental health in overall well-being...',
    image: 'https://picsum.photos/seed/health3/800/400'
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
</style> 