<template>
  <div>
    <!-- Contact Header -->
    <section class="contact-header">
      <v-container class="py-16">
        <v-row justify="center">
          <v-col cols="12" md="8" class="text-center">
            <h1 class="text-h2 font-weight-bold mb-6">
              Contact Us
            </h1>
            <p class="text-h6">
              Get in touch with our healthcare professionals
            </p>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- Contact Form -->
    <section class="contact-form">
      <v-container class="py-16">
        <v-row justify="center">
          <v-col cols="12" md="8">
            <v-card class="pa-6">
              <v-form @submit.prevent="submitForm">
                <v-row>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.name"
                      label="Full Name"
                      required
                    ></v-text-field>
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.email"
                      label="Email"
                      type="email"
                      required
                    ></v-text-field>
                  </v-col>
                  <v-col cols="12">
                    <v-text-field
                      v-model="form.subject"
                      label="Subject"
                      required
                    ></v-text-field>
                  </v-col>
                  <v-col cols="12">
                    <v-textarea
                      v-model="form.message"
                      label="Message"
                      required
                    ></v-textarea>
                  </v-col>
                  <v-col cols="12" class="text-center">
                    <v-btn
                      color="primary"
                      size="large"
                      type="submit"
                      rounded
                    >
                      Send Message
                    </v-btn>
                  </v-col>
                </v-row>
              </v-form>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- Contact Information -->
    <section class="contact-info bg-grey-lighten-4">
      <v-container class="py-16">
        <v-row>
          <v-col cols="12" md="4">
            <div class="text-center">
              <v-icon size="48" color="primary" class="mb-4">mdi-map-marker</v-icon>
              <h3 class="text-h6 mb-2">Address</h3>
              <p>123 Healthcare Ave, Medical City</p>
            </div>
          </v-col>
          <v-col cols="12" md="4">
            <div class="text-center">
              <v-icon size="48" color="primary" class="mb-4">mdi-phone</v-icon>
              <h3 class="text-h6 mb-2">Phone</h3>
              <p>(555) 123-4567</p>
            </div>
          </v-col>
          <v-col cols="12" md="4">
            <div class="text-center">
              <v-icon size="48" color="primary" class="mb-4">mdi-email</v-icon>
              <h3 class="text-h6 mb-2">Email</h3>
              <p>info@healis.com</p>
            </div>
          </v-col>
        </v-row>
      </v-container>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const form = ref({
  name: '',
  email: '',
  subject: '',
  message: ''
})

const submitForm = () => {
  // Handle form submission
  console.log('Form submitted:', form.value)
  // Reset form
  form.value = {
    name: '',
    email: '',
    subject: '',
    message: ''
  }
}
</script>

<style scoped>
.contact-header {
  background: linear-gradient(135deg, var(--v-primary-base) 0%, var(--v-primary-darken2) 100%);
  color: white;
}

.contact-info {
  background-color: #f8f9fa;
}
</style> 