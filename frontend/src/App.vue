<template>
  <div>
    <Header />
    <CommitLog class="pt-32 sm:pt-40" :commits="commits" />
  </div>
</template>

<script>
import { reactive } from 'vue';

import Header from '@/components/Header';
import CommitLog from '@/components/CommitLog';

export default {
  components: {
    Header,
    CommitLog,
  },
  setup() {
    const commits = reactive([]);

    fetch('/api/commitlog')
      .then((res) => res.json())
      .then((res) => {
        commits.push(...res.log);
      });

    return { commits };
  },
};
</script>
