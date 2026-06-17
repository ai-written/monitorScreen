<template>
  <div class="metric-card">
    <div class="metric-label">{{ label }}</div>
    <div class="metric-value">
      <span class="metric-num">{{ formattedValue }}</span>
      <span class="metric-unit">{{ unit }}</span>
    </div>
    <div class="progress-bar">
      <div class="progress-fill" :style="fillStyle"></div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    label: { type: String, required: true },
    value: { type: Number, default: 0 },
    unit: { type: String, default: '' },
    max: { type: Number, default: 100 },
    barColor: { type: String, default: 'var(--cyan)' },
    decimals: { type: Number, default: 1 },
    showZero: { type: Boolean, default: true },
  },
  computed: {
    formattedValue() {
      if (!this.showZero && this.value === 0) return '--'
      return this.value.toFixed(this.decimals)
    },
    pct() {
      if (this.max <= 0) return 0
      const p = (this.value / this.max) * 100
      return Math.min(p, 100)
    },
    fillStyle() {
      return {
        width: this.pct + '%',
        background: this.barColor
      }
    }
  }
}
</script>

<style scoped>
.metric-card {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 8px;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
}

.metric-label {
  font-size: 9px;
  color: var(--text-dim);
  letter-spacing: 1.5px;
  margin-bottom: 4px;
  text-transform: uppercase;
}

.metric-value {
  display: flex;
  align-items: baseline;
  gap: 4px;
  margin-bottom: 8px;
}

.metric-num {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
}

.metric-unit {
  font-size: 13px;
  color: var(--text-secondary);
}

.progress-bar {
  height: 4px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.4s ease;
}
</style>
