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
      commits: [
        {
          id: '123123123',
          created_at: '2018/04/02 03:04 PM',
          author: 'SpinyNuts',
          message: 'Deeeeeze nuts for fucks sake',
          avatar_url: 'https://avatars.githubusercontent.com/u/10344453?',
          link:
            'https://github.com/BenMusch/mittab-deploy/commit/d22eef887c6966db2d4f85a5942a1575151482c3',
        },
        {
          id: '1231231356',
          created_at: '2018/04/02 01:21 PM',
          author: 'joe.frazer',
          message: 'Add a fucking sleep because fuck it',
          avatar_url: 'https://avatars.githubusercontent.com/u/2784192?',
          link:
            'https://github.com/rebel1324/CityRP/commit/06b2ecab116a6e86ac708b8ef0cb8f5fe95b8adf',
        },
      ],
    };
  },
};
</script>
