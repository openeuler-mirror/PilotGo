export let pickerOptions = {
  shortcuts: [{
    text: '最近10分钟',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 10 * 60 * 1000);
      picker.$emit('pick', [start, end]);
    }
  }, 
  {
    text: '最近30分钟',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 30 * 60 * 1000);
      picker.$emit('pick', [start, end]);
    }
  },
  {
    text: '最近2小时',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 2 * 60 * 60 * 1000);
      picker.$emit('pick', [start, end]);
    }
  },
  {
    text: '最近6小时',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 6 * 60 * 60 * 1000);
      picker.$emit('pick', [start, end]);
    }
  },
  {
    text: '最近12小时',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 12 * 60 * 60 * 1000);
      picker.$emit('pick', [start, end]);
    }
  },
  {
    text: '最近24小时',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 24 * 60 * 60 * 1000);
      picker.$emit('pick', [start, end]);
    }
  },{
    text: '最近3天',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 3);
      picker.$emit('pick', [start, end]);
    }
  },{
    text: '最近1周',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
      picker.$emit('pick', [start, end]);
    }
  },{
    text: '最近1个月',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
      picker.$emit('pick', [start, end]);
    }
  }, {
    text: '最近3个月',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30 * 3);
      picker.$emit('pick', [start, end]);
    }
  }, {
    text: '最近6个月',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30 * 6);
      picker.$emit('pick', [start, end]);
    }
  }, {
    text: '最近1年',
    onClick(picker) {
      const end = new Date();
      const start = new Date();
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30 *12);
      picker.$emit('pick', [start, end]);
    }
  }]
}

