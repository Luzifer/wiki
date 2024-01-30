<template>
  <div class="container">
    <div class="row">
      <div
        ref="content"
        class="col relAnchor"
      >
        <router-link
          v-if="$route.params.page"
          class="btn btn-secondary btn-sm editBtn"
          :to="{ name: 'edit', params: { page: $route.params.page } }"
        >
          <i class="fas fa-edit" />
        </router-link>

        <md-render
          :content="content"
          :prerender="prerender"
          @rendered="rendered"
        />
      </div>
    </div>
  </div>
</template>

<script>
/* global Prism */

import mdRender from './markdown.vue'

export default {

  components: {
    mdRender,
  },

  data() {
    return {
      content: '',
    }
  },

  methods: {
    intLinkClick(evt) {
      const link = evt.target
      this.$router.push({ name: 'view', params: { page: link.dataset.page } })
      return false
    },

    loadPage(pageName) {
      console.debug(`Loading ${pageName}...`)
      return fetch(`/_content/${pageName}`)
        .then(resp => {
          if (resp.status === 404) {
            this.$router.push({ name: 'edit', params: { page: pageName } })
            return
          }

          return resp.json()
        })
        .then(data => {
          if (!data) {
            return
          }

          this.content = data.content
        })
        .catch(err => {
          console.error(err)
        })
    },

    prerender(mdtext) {
      // replace [[Internal]] links
      mdtext = mdtext.replace(
        new RegExp(/\[\[([^\]]+)\]\]/, 'g'),
        '<a class="intLink" data-page="$1" href="$1">$1</a>',
      )

      return mdtext
    },

    rendered() {
      // Give the DOM a moment to update before manipulating further
      window.setTimeout(() => {
        // Add listeners to internal links
        const links = this.$refs.content.getElementsByClassName('intLink')
        for (const link of links) {
          link.onclick = this.intLinkClick
        }

        // Highlight code blocks
        Prism.highlightAll()
      }, 100)
    },
  },

  mounted() {
    this.loadPage(this.$route.params.page)
  },

  name: 'WikiView',

  watch: {
    '$route'(to, from) {
      if (to.params.page === from.params.page) {
        return
      }

      this.loadPage(to.params.page)
    },
  },
}
</script>

<style>
.editBtn {
  position: absolute;
  right: 5px;
  top: 5px;
}
.relAnchor {
  position: relative;
}
</style>
