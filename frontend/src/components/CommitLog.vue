<template>
  <div class="m-auto lg:max-w-4xl px-4 lg:px-0">
    <ul id="commits-log">
      <CommitLogItem
        v-for="commit in commits"
        :key="commit.id"
        v-bind="commit"
      />
      <Spinner v-if="busy" class="mt-10 mx-auto h-16 w-16" />
    </ul>
  </div>
</template>

<script>
import CommitLogItem from './CommitLogItem';
import Spinner from './Spinner';

export default {
  components: {
    CommitLogItem,
    Spinner,
  },
  props: {
    commits: {
      type: Array,
      default: () => [],
    },
    busy: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
