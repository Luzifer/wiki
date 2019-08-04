<template>
  <b-container>
    <b-row>
      <b-col
        ref="content"
        class="relAnchor"
      >
        <b-btn
          class="editBtn"
          variant="secondary"
          size="sm"
          :to="{ name: 'edit', params: { page: $route.params.page } }"
        >
          <i class="fas fa-edit" />
        </b-btn>
        <vue-markdown
          :source="content"
          :prerender="prerender"
          @rendered="rendered"
        />
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import axios from 'axios'
import VueMarkdown from 'vue-markdown'

export default {
  name: 'View',

  components: {
    VueMarkdown,
  },

  data() {
    return {
      content: '',
    }
  },

  watch: {
    '$route'(to, from) {
      if (to.params.page === from.params.page) {
        return
      }

      this.loadPage(to.params.page)
    },
  },

  mounted() {
    this.loadPage(this.$route.params.page)
  },

  methods: {
    intLinkClick(evt) {
      const link = evt.target
      this.$router.push({ name: 'view', params: { page: link.dataset.page } })
      return false
    },

    loadPage(pageName) {
      console.debug(`Loading ${pageName}...`)
      axios.get(`/_content/${pageName}`)
        .then(resp => {
          this.content = resp.data.content
        })
        .catch(err => {
          if (err.response && err.response.status === 404) {
            this.$router.push({ name: 'edit', params: { page: pageName } })
            return
          }
          console.error(err)
          // FIXME: Show error
        })
    },

    prerender(mdtext) {
      // replace [[Internal]] links
      mdtext = mdtext.replace(/\[\[([^\]]+)\]\]/, '<a class="intLink" data-page="$1" href="$1">$1</a>')

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
