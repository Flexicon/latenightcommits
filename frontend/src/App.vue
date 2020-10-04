<template>
  <div>
    <Header />
    <CommitLog :commits="commits" />
  </div>
</template>

<script>
import Header from '@/components/Header';
import CommitLog from '@/components/CommitLog';

export default {
  components: {
    Header,
    CommitLog,
  },
  data() {
    return {
      commits: [],
    };
  },
  async created() {
    const res = await fetch('/api/commitslog').then((res) => res.json());
    this.commits = res.log;
  },
};
</script>
