<template>
  <div class="m-auto lg:max-w-4xl px-4 lg:px-0">
    <ul id="commits-log" class="divide-y divide-gray-400">
      <CommitLogItem
        v-for="commit in commits"
        :key="commit.id"
        v-bind="commit"
      />
    </ul>
  </div>
</template>

<script>
import CommitLogItem from './CommitLogItem';

export default {
  components: {
    CommitLogItem,
  },
  props: {
    commits: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
