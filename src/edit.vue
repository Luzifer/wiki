<template>
  <div class="container">
    <div class="row">
      <div class="col">
        <textarea ref="editor" />

        <button
          class="btn btn-primary"
          @click="save"
        >
          Save
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import EasyMDE from 'easymde'

export default {
  data() {
    return {
      editor: null,
    }
  },

  methods: {
    loadPage(pageName) {
      console.debug(`Loading ${pageName}...`)
      return fetch(`/_content/${pageName}`)
        .then(resp => {
          if (resp.status === 404) {
            return { content: `# ${pageName}` }
          }

          return resp.json()
        })
        .then(data => {
          if (this.editor) {
            this.editor.toTextArea()
            this.editor = null
          }

          this.editor = new EasyMDE({
            element: this.$refs.editor,
            forceSync: true,
            indentWithTabs: false,
            initialValue: data.content,
          })

          // this.editor.codemirror.setValue(data.content)
        })
        .catch(err => {
          if (err.response && err.response.status === 404) {
            return
          }
          console.error(err)
        })
    },

    save() {
      return fetch(`/_content/${this.$route.params.page}`, {
        body: JSON.stringify({ content: this.$refs.editor.value }),
        method: 'POST',
      })
        .then(() => {
          this.$router.push({ name: 'view', params: { page: this.$route.params.page } })
        })
        .catch(err => {
          console.error(err)
        })
    },
  },

  mounted() {
    this.loadPage(this.$route.params.page)
  },

  name: 'WikiEdit',

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
