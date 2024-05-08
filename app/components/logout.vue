<template>
  <div class="no-content">
    <p>Logging out...</p>
  </div>
</template>
<script>
export default {
  created() {
    window.utils.title('Log Out');
    this.$http.post('/app/logout').then(function() {
      location.href = '/';
    })
  }
}
</script>
