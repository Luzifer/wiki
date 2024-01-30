<template>
  <div>
    <nav class="navbar navbar-expand-lg bg-body-tertiary mb-3">
      <div class="container-fluid">
        <router-link
          class="navbar-brand"
          :to="{ name: 'home' }"
        >
          Wiki
        </router-link>
        <button
          class="navbar-toggler"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#navbarSupportedContent"
          aria-controls="navbarSupportedContent"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span class="navbar-toggler-icon" />
        </button>
        <div
          id="navbarSupportedContent"
          class="collapse navbar-collapse"
        >
          <ul class="navbar-nav me-auto mb-2 mb-lg-0">
            <li
              v-for="page in navContent"
              :key="page"
              class="nav-item"
            >
              <router-link
                class="nav-link"
                :to="{ name: 'view', params: { page } }"
              >
                {{ page }}
              </router-link>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    <router-view />
  </div>
</template>

<script>
export default {
  data() {
    return {
      navContent: '',
    }
  },

  methods: {
    loadNav() {
      return fetch(`/_content/_navigation`)
        .then(resp => resp.json())
        .then(data => {
          this.navContent = data.content.split('\n')
            .filter(el => el && !el.match(/^!/))
        })
        .catch(err => {
          console.error(err)
        })
    },
  },

  mounted() {
    this.loadNav()
  },

  name: 'WikiApp',
}
</script>
