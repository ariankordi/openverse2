<template></template>
<script>
export default {
  beforeMount() {
    alert('There\'s no ' + this.$route.path.substring(1) + ' on this site, so we\'ll just redirect you home.');
    //this.router.go(-1);
    this.$router.push('/');
  }
}
</script>
