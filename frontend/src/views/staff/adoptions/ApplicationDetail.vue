<template>
  <div class="application-detail">
    <LoadingSpinner v-if="loading" />

    <div
      v-else-if="application"
      class="detail-container"
    >
      <div class="detail-header">
        <Button
          icon="pi pi-arrow-left"
          class="p-button-text"
          @click="router.back()"
        />
        <h1>{{ $t('adoption.applicationDetail') }}</h1>
        <div
          v-if="canReview"
          class="header-actions"
        >
          <Button
            :label="$t('adoption.approve')"
            icon="pi pi-check"
            class="p-button-success"
            @click="showApproveDialog"
          />
          <Button
            :label="$t('adoption.reject')"
            icon="pi pi-times"
            class="p-button-danger"
            @click="showRejectDialog"
          />
        </div>
      </div>

      <Card class="status-card">
        <template #content>
          <div class="status-info">
            <div class="status-item">
              <label>{{ $t('adoption.status') }}</label>
              <Badge :variant="getStatusVariant(application.status)">
                {{ $t(`adoption.${application.status}`) }}
              </Badge>
            </div>
            <div class="status-item">
              <label>{{ $t('adoption.applicationDate') }}</label>
              <p>{{ formatDate(application.application_date) }}</p>
            </div>
            <div
              v-if="application.animal"
              class="status-item"
            >
              <label>{{ $t('animal.name') }}</label>
              <router-link
                :to="`/staff/animals/${application.animal.id}`"
                class="animal-link"
              >
                {{ application.animal.name }}
              </router-link>
            </div>
          </div>
        </template>
      </Card>

      <TabView>
        <TabPanel :header="$t('adoption.applicantInfo')">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('adoption.firstName') }}</label>
                  <p>{{ application.applicant_first_name }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('adoption.lastName') }}</label>
                  <p>{{ application.applicant_last_name }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('adoption.email') }}</label>
                  <p>{{ application.email }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('adoption.phone') }}</label>
                  <p>{{ application.phone || 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('adoption.dateOfBirth') }}</label>
                  <p>{{ formatDate(application.date_of_birth) }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('adoption.occupation') }}</label>
                  <p>{{ application.occupation || 'N/A' }}</p>
                </div>
              </div>
              <div
                v-if="application.address"
                class="info-item full-width"
              >
                <label>{{ $t('adoption.address') }}</label>
                <p>
                  {{ application.address.street }}<br>
                  {{ application.address.city }}, {{ application.address.state }} {{ application.address.postal_code }}<br>
                  {{ application.address.country }}
                </p>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel :header="$t('adoption.householdInfo')">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('adoption.householdType') }}</label>
                  <p>{{ application.household_type }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('adoption.hasYard') }}</label>
                  <Badge :variant="application.has_yard ? 'success' : 'neutral'">
                    {{ application.has_yard ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('adoption.yardFenced') }}</label>
                  <Badge :variant="application.yard_fenced ? 'success' : 'neutral'">
                    {{ application.yard_fenced ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('adoption.householdMembers') }}</label>
                  <p>{{ application.household_members || 'N/A' }}</p>
                </div>
              </div>

              <div
                v-if="application.other_pets && application.other_pets.length > 0"
                class="info-item full-width"
              >
                <label>{{ $t('adoption.otherPets') }}</label>
                <div class="pets-list">
                  <Card
                    v-for="(pet, index) in application.other_pets"
                    :key="index"
                    class="pet-card"
                  >
                    <template #content>
                      <p><strong>{{ pet.type }}</strong> - {{ pet.breed || 'Mixed' }}</p>
                      <p>{{ $t('animal.age') }}: {{ pet.age }} {{ $t('home.adoptions.years') }}</p>
                      <p>{{ $t('animal.spayedNeutered') }}: {{ pet.spayed_neutered ? $t('common.yes') : $t('common.no') }}</p>
                    </template>
                  </Card>
                </div>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel :header="$t('adoption.experience')">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('adoption.yearsOfExperience') }}</label>
                  <p>{{ application.years_of_experience || 0 }} {{ $t('home.adoptions.years') }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('adoption.previousAdoptions') }}</label>
                  <p>{{ application.previous_adoptions || 'N/A' }}</p>
                </div>
              </div>

              <div
                v-if="application.veterinarian_info"
                class="info-item full-width"
              >
                <label>{{ $t('adoption.veterinarianInfo') }}</label>
                <p>
                  <strong>{{ application.veterinarian_info.name }}</strong><br>
                  {{ application.veterinarian_info.clinic }}<br>
                  {{ $t('adoption.phone') }}: {{ application.veterinarian_info.phone }}
                </p>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel :header="$t('adoption.references')">
          <Card>
            <template #content>
              <div
                v-if="application.references && application.references.length > 0"
                class="references-list"
              >
                <Card
                  v-for="(reference, index) in application.references"
                  :key="index"
                  class="reference-card"
                >
                  <template #content>
                    <div class="reference-info">
                      <h3>{{ reference.name }}</h3>
                      <p><strong>{{ $t('adoption.references') }}:</strong> {{ reference.relationship }}</p>
                      <p><strong>{{ $t('adoption.phone') }}:</strong> {{ reference.phone }}</p>
                      <p><strong>{{ $t('adoption.yearsOfExperience') }}:</strong> {{ reference.years_known }} {{ $t('home.adoptions.years') }}</p>
                    </div>
                  </template>
                </Card>
              </div>
              <EmptyState
                v-else
                :message="$t('adoption.noApplicationsFound')"
              />
            </template>
          </Card>
        </TabPanel>

        <TabPanel :header="$t('adoption.employment')">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('adoption.employmentStatus') }}</label>
                  <p>{{ application.employment_status || 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('adoption.workHours') }}</label>
                  <p>{{ application.work_hours_per_day || 0 }}h / {{ $t('common.day') }}</p>
                </div>
                <div class="info-item full-width">
                  <label>{{ $t('adoption.whoWillCare') }}</label>
                  <p>{{ application.who_will_care || 'N/A' }}</p>
                </div>
              </div>

              <div class="info-item full-width">
                <label>{{ $t('adoption.reasonForAdoption') }}</label>
                <p>{{ application.reason_for_adoption || 'N/A' }}</p>
              </div>

              <div class="info-item full-width">
                <label>{{ $t('adoption.expectations') }}</label>
                <p>{{ application.expectations || 'N/A' }}</p>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel
          v-if="application.review_notes || application.reviewed_by"
          :header="$t('adoption.reviewNotes')"
        >
          <Card>
            <template #content>
              <div class="info-grid">
                <div
                  v-if="application.reviewed_by"
                  class="info-item"
                >
                  <label>Reviewed By</label>
                  <p>{{ application.reviewed_by }}</p>
                </div>
                <div
                  v-if="application.reviewed_at"
                  class="info-item"
                >
                  <label>Reviewed At</label>
                  <p>{{ formatDate(application.reviewed_at) }}</p>
                </div>
              </div>
              <div
                v-if="application.review_notes"
                class="info-item full-width"
              >
                <label>{{ $t('adoption.reviewNotes') }}</label>
                <p>{{ application.review_notes }}</p>
              </div>
            </template>
          </Card>
        </TabPanel>
      </TabView>
    </div>

    <!-- Approve Dialog -->
    <Dialog
      v-model:visible="approveDialogVisible"
      :header="$t('adoption.approve')"
      :modal="true"
      style="width: 30rem"
    >
      <div class="dialog-content">
        <label for="approve-notes">{{ $t('adoption.reviewNotes') }}</label>
        <Textarea
          id="approve-notes"
          v-model="reviewNotes"
          rows="4"
          :placeholder="$t('adoption.reviewNotes')"
        />
      </div>
      <template #footer>
        <Button
          :label="$t('common.cancel')"
          class="p-button-secondary"
          @click="approveDialogVisible = false"
        />
        <Button
          :label="$t('adoption.approve')"
          class="p-button-success"
          :loading="submitting"
          @click="handleApprove"
        />
      </template>
    </Dialog>

    <!-- Reject Dialog -->
    <Dialog
      v-model:visible="rejectDialogVisible"
      :header="$t('adoption.reject')"
      :modal="true"
      style="width: 30rem"
    >
      <div class="dialog-content">
        <label for="reject-reason">Reason *</label>
        <Textarea
          id="reject-reason"
          v-model="rejectReason"
          rows="4"
          placeholder="Please provide a reason for rejection"
          required
        />
      </div>
      <template #footer>
        <Button
          :label="$t('common.cancel')"
          class="p-button-secondary"
          @click="rejectDialogVisible = false"
        />
        <Button
          :label="$t('adoption.reject')"
          class="p-button-danger"
          :loading="submitting"
          :disabled="!rejectReason"
          @click="handleReject"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { adoptionService } from '@/services/adoptionService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import Dialog from 'primevue/dialog'
import Textarea from 'primevue/textarea'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const application = ref(null)
const loading = ref(true)
const submitting = ref(false)
const approveDialogVisible = ref(false)
const rejectDialogVisible = ref(false)
const reviewNotes = ref('')
const rejectReason = ref('')

const canReview = computed(() => {
  return application.value &&
    (application.value.status === 'pending' || application.value.status === 'under_review')
})

const loadApplication = async () => {
  try {
    loading.value = true
    application.value = await adoptionService.getApplication(route.params.id)
  } catch (error) {
    console.error('Error loading application:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load application',
      life: 3000
    })
    router.push('/staff/adoptions/applications')
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString()
}

const getStatusVariant = (status) => ({
  pending: 'warning',
  under_review: 'info',
  approved: 'success',
  rejected: 'danger',
  withdrawn: 'neutral'
}[status] || 'neutral')

const showApproveDialog = () => {
  reviewNotes.value = ''
  approveDialogVisible.value = true
}

const showRejectDialog = () => {
  rejectReason.value = ''
  rejectDialogVisible.value = true
}

const handleApprove = async () => {
  try {
    submitting.value = true
    await adoptionService.approveApplication(application.value.id, reviewNotes.value)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: t('adoption.applicationApproved'),
      life: 3000
    })
    approveDialogVisible.value = false
    await loadApplication()
  } catch (error) {
    console.error('Error approving application:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to approve application',
      life: 3000
    })
  } finally {
    submitting.value = false
  }
}

const handleReject = async () => {
  if (!rejectReason.value) {
    toast.add({
      severity: 'warn',
      summary: 'Warning',
      detail: 'Please provide a reason for rejection',
      life: 3000
    })
    return
  }

  try {
    submitting.value = true
    await adoptionService.rejectApplication(application.value.id, rejectReason.value)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: t('adoption.applicationRejected'),
      life: 3000
    })
    rejectDialogVisible.value = false
    await loadApplication()
  } catch (error) {
    console.error('Error rejecting application:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to reject application',
      life: 3000
    })
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadApplication()
})
</script>

<style scoped>
.application-detail {
  max-width: 1200px;
  margin: 0 auto;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
}

.detail-header h1 {
  flex: 1;
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.status-card {
  margin-bottom: 1.5rem;
}

.status-info {
  display: flex;
  gap: 2rem;
  align-items: center;
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.status-item label {
  font-weight: 600;
  color: #6b7280;
  font-size: 0.875rem;
  text-transform: uppercase;
}

.status-item p {
  color: #2c3e50;
  font-size: 1rem;
  margin: 0;
}

.animal-link {
  color: #3b82f6;
  text-decoration: none;
  font-weight: 600;
  font-size: 1rem;
}

.animal-link:hover {
  text-decoration: underline;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-item label {
  font-weight: 600;
  color: #6b7280;
  font-size: 0.875rem;
  text-transform: uppercase;
}

.info-item p {
  color: #2c3e50;
  font-size: 1rem;
  margin: 0;
  line-height: 1.6;
}

.full-width {
  grid-column: 1 / -1;
}

.pets-list, .references-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
  margin-top: 1rem;
}

.pet-card, .reference-card {
  border: 1px solid #e5e7eb;
}

.reference-info h3 {
  margin: 0 0 0.5rem 0;
  color: #2c3e50;
  font-size: 1.125rem;
}

.reference-info p {
  margin: 0.25rem 0;
  font-size: 0.875rem;
}

.dialog-content {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.dialog-content label {
  font-weight: 600;
  color: #374151;
}
</style>
