<template>
  <div>
    <Header />
    <CommitLog :commits="commits" :busy="busy" />
  </div>
</template>

<script>
import { reactive, ref } from 'vue';

import Header from '@/components/Header.vue';
import CommitLog from '@/components/CommitLog.vue';

const API_BASE_URL = (
  import.meta.env.VITE_API_BASE_URL ||
  'https://latenightcommits-api.nerfthis.xyz'
).replace(/\/+$/, '');

export default {
  components: {
    Header,
    CommitLog,
  },
  setup() {
    const commits = reactive([]);
    let page = 1;
    let busy = ref(false);
    let hasNextPage = true;

    const fetchNextLogPage = () => {
      busy.value = true;
      fetch(`${API_BASE_URL}/commitlog?page=${page}&per_page=50`)
        .then((res) => res.json())
        .then((res) => {
          commits.push(...res.log);
          hasNextPage = Boolean(res.has_next_page);
          page += 1;
        })
        .finally(() => {
          busy.value = false;
        });
    };

    window.addEventListener('scroll', () => {
      if (!busy.value && hasNextPage) {
        const doc = document.documentElement;
        const currentPosition = doc.scrollTop + window.innerHeight;

        if (doc.scrollHeight - currentPosition < window.innerHeight) {
          fetchNextLogPage();
        }
      }
    });
    fetchNextLogPage();

    return { busy, commits };
  },
};
</script>

<style>
body, html {
  font-family: 'Roboto', sans-serif;
  background-color: #374d78;
}
</style>
