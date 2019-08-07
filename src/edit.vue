<template>
  <b-container>
    <b-row>
      <b-col>
        <textarea ref="editor" />

        <b-btn
          variant="primary"
          @click="save"
        >
          Save
        </b-btn>
      </b-col>
    </b-row>
  </b-container>
</template>

<script>
import axios from 'axios'
import EasyMDE from 'easymde'

export default {
  name: 'Edit',

  data() {
    return {
      editor: null,
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
    if (!this.editor) {
      this.editor = new EasyMDE({
        element: this.$refs.editor,
        forceSync: true,
        indentWithTabs: false,
      })

      window.editor = this.editor
    }
    this.loadPage(this.$route.params.page)
  },

  methods: {
    loadPage(pageName) {
      console.debug(`Loading ${pageName}...`)
      axios.get(`/_content/${pageName}`)
        .then(resp => {
          this.editor.codemirror.setValue(resp.data.content)
        })
        .catch(err => {
          if (err.response && err.response.status === 404) {
            return
          }
          console.error(err)
          // FIXME: Show error
        })
    },

    save() {
      axios.post(`/_content/${this.$route.params.page}`, {
        content: this.$refs.editor.value,
      })
        .then(() => {
          this.$router.push({ name: 'view', params: { page: this.$route.params.page } })
        })
        .catch(err => {
          console.error(err)
        })
    },
  },
}
</script>

<style>
.editor-toolbar {
  background-color: #ccc;
}
</style>
