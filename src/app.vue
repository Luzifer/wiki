<template>
  <div>
    <b-navbar
      type="dark"
      variant="primary"
      class="mb-4"
    >
      <b-navbar-brand
        href="/"
        @click.prevent="$router.push({ name: 'home' })"
      >
        Wiki
      </b-navbar-brand>

      <b-navbar-nav>
        <b-nav-item
          v-for="page in navContent"
          :key="page"
          :to="{ name: 'view', params: { page } }"
        >
          {{ page }}
        </b-nav-item>
      </b-navbar-nav>

      <template ref="nav">
        <vue-markdown
          class="nav"
          :source="navContent"
          :prerender="prerender"
          @rendered="rendered"
        />
      </template>
    </b-navbar>

    <router-view />
  </div>
</template>

<script>
import axios from 'axios'
import VueMarkdown from 'vue-markdown'

export default {
  name: 'App',

  components: {
    VueMarkdown,
  },

  data() {
    return {
      navContent: '',
    }
  },

  mounted() {
    this.loadNav()
  },

  methods: {
    loadNav() {
      axios.get(`/_content/_navigation`)
        .then(resp => {
          this.navContent = resp.data.content.split('\n')
            .filter(el => !el.match(/^!/))
        })
        .catch(err => {
          if (err.response && err.response.status === 404) {
            return
          }
          console.error(err)
        })
    },
  },
}
</script>
