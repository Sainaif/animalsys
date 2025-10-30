<template>
  <div class="campaign-form-page">
    <div class="page-header">
      <h1 class="page-title">
        {{ isEdit ? t('campaigns.editCampaign') : t('campaigns.addCampaign') }}
      </h1>
    </div>

    <BaseCard>
      <LoadingSpinner v-if="loading" />
      <form v-else @submit.prevent="handleSubmit">
        <div class="form-section">
          <h3 class="section-title">{{ t('campaigns.basicInfo') }}</h3>

          <div class="form-row">
            <FormGroup :label="t('campaigns.name')" :error="errors.name" required>
              <input
                v-model="form.name"
                type="text"
                class="form-control"
                :class="{ 'error': errors.name }"
                :placeholder="t('campaigns.namePlaceholder')"
              />
            </FormGroup>

            <FormGroup :label="t('campaigns.campaignType')" :error="errors.type" required>
              <select v-model="form.type" class="form-control" :class="{ 'error': errors.type }">
                <option value="">{{ t('campaigns.selectType') }}</option>
                <option value="fundraising">{{ t('campaigns.typeFundraising') }}</option>
                <option value="adoption">{{ t('campaigns.typeAdoption') }}</option>
                <option value="event">{{ t('campaigns.typeEvent') }}</option>
                <option value="awareness">{{ t('campaigns.typeAwareness') }}</option>
              </select>
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('common.description')" :error="errors.description">
              <textarea
                v-model="form.description"
                class="form-control"
                :placeholder="t('campaigns.descriptionPlaceholder')"
                rows="4"
              ></textarea>
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('campaigns.startDate')" :error="errors.start_date" required>
              <input
                v-model="form.start_date"
                type="date"
                class="form-control"
                :class="{ 'error': errors.start_date }"
              />
            </FormGroup>

            <FormGroup :label="t('campaigns.endDate')" :error="errors.end_date" required>
              <input
                v-model="form.end_date"
                type="date"
                class="form-control"
                :class="{ 'error': errors.end_date }"
              />
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('common.status')" :error="errors.status">
              <select v-model="form.status" class="form-control">
                <option value="active">{{ t('campaigns.statusActive') }}</option>
                <option value="upcoming">{{ t('campaigns.statusUpcoming') }}</option>
                <option value="completed">{{ t('campaigns.statusCompleted') }}</option>
                <option value="cancelled">{{ t('campaigns.statusCancelled') }}</option>
              </select>
            </FormGroup>
          </div>
        </div>

        <div v-if="form.type === 'fundraising'" class="form-section">
          <h3 class="section-title">{{ t('campaigns.goalInfo') }}</h3>
          <div class="form-row">
            <FormGroup :label="t('campaigns.goalAmount')" :error="errors.goal_amount">
              <input
                v-model.number="form.goal_amount"
                type="number"
                step="0.01"
                min="0"
                class="form-control"
                :placeholder="t('campaigns.goalAmountPlaceholder')"
              />
            </FormGroup>
          </div>
        </div>

        <div v-if="form.type === 'adoption'" class="form-section">
          <h3 class="section-title">{{ t('campaigns.goalInfo') }}</h3>
          <div class="form-row">
            <FormGroup :label="t('campaigns.goalAdoptions')" :error="errors.goal_adoptions">
              <input
                v-model.number="form.goal_adoptions"
                type="number"
                min="0"
                class="form-control"
                :placeholder="t('campaigns.goalAdoptionsPlaceholder')"
              />
            </FormGroup>
          </div>
        </div>

        <div class="form-actions">
          <BaseButton type="button" variant="secondary" @click="goBack">
            {{ t('common.cancel') }}
          </BaseButton>
          <BaseButton type="submit" variant="primary" :disabled="submitting">
            {{ submitting ? t('common.saving') : t('common.save') }}
          </BaseButton>
        </div>
      </form>
    </BaseCard>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useNotificationStore } from '../../stores/notifications'
import { API } from '../../api'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import FormGroup from '../../components/base/FormGroup.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const notificationStore = useNotificationStore()

const isEdit = ref(false)
const loading = ref(false)
const submitting = ref(false)

const form = reactive({
  name: '',
  type: '',
  description: '',
  start_date: '',
  end_date: '',
  status: 'upcoming',
  goal_amount: null,
  goal_adoptions: null,
})

const errors = reactive({})

async function fetchCampaign() {
  if (!route.params.id) return

  try {
    loading.value = true
    const response = await API.campaigns.getById(route.params.id)
    const campaign = response.data

    Object.keys(form).forEach(key => {
      if (campaign[key] !== undefined) {
        form[key] = campaign[key]
      }
    })

    if (form.start_date) {
      form.start_date = new Date(form.start_date).toISOString().split('T')[0]
    }
    if (form.end_date) {
      form.end_date = new Date(form.end_date).toISOString().split('T')[0]
    }
  } catch (error) {
    console.error('Failed to fetch campaign:', error)
    notificationStore.error(t('campaigns.fetchError'))
    goBack()
  } finally {
    loading.value = false
  }
}

function validateForm() {
  Object.keys(errors).forEach(key => delete errors[key])
  let isValid = true

  if (!form.name || form.name.trim().length === 0) {
    errors.name = t('common.required')
    isValid = false
  }

  if (!form.type) {
    errors.type = t('common.required')
    isValid = false
  }

  if (!form.start_date) {
    errors.start_date = t('common.required')
    isValid = false
  }

  if (!form.end_date) {
    errors.end_date = t('common.required')
    isValid = false
  }

  if (form.start_date && form.end_date && new Date(form.end_date) < new Date(form.start_date)) {
    errors.end_date = t('campaigns.endDateError')
    isValid = false
  }

  return isValid
}

async function handleSubmit() {
  if (!validateForm()) {
    notificationStore.error(t('common.fixErrors'))
    return
  }

  try {
    submitting.value = true

    if (isEdit.value) {
      await API.campaigns.update(route.params.id, form)
      notificationStore.success(t('campaigns.updateSuccess'))
    } else {
      await API.campaigns.create(form)
      notificationStore.success(t('campaigns.createSuccess'))
    }

    goBack()
  } catch (error) {
    console.error('Failed to save campaign:', error)
    notificationStore.error(
      isEdit.value ? t('campaigns.updateError') : t('campaigns.createError')
    )
  } finally {
    submitting.value = false
  }
}

function goBack() {
  router.push({ name: 'Campaigns' })
}

onMounted(() => {
  if (route.params.id) {
    isEdit.value = true
    fetchCampaign()
  }
})
</script>

<style scoped>
.campaign-form-page {
  max-width: 900px;
  padding: 2rem;
}

.page-header {
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  margin: 0;
}

.form-section {
  margin-bottom: 2rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--border-color);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-control {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-control:focus {
  outline: none;
  border-color: var(--primary-color);
}

.form-control.error {
  border-color: var(--danger-color);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
}
</style>
