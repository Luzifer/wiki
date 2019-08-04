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
    </b-navbar>

    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App',
}
</script>
