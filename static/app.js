// ===== Fingo - Frontend Application =====

const API = {
  users: "/users",
  transactions: "/transactions",
  goals: "/goals",
};

// ========================================================================
//  I18N — INTERNATIONALIZATION
// ========================================================================

const LANG_KEY = "fingo-lang";

const translations = {
  en: {
    // Header
    subtitle: "Personal Financial Manager",

    // Tabs
    tab_users: "Users",
    tab_transactions: "Transactions",
    tab_goals: "Goals",

    // Users
    users_title: "Users",
    btn_new_user: "+ New User",
    form_create_user: "Create User",
    form_edit_user: "Edit User",
    label_username: "Username",
    ph_username: "Enter username",
    label_current_amount: "Current Amount ($)",
    label_monthly_inputs: "Monthly Inputs ($)",
    label_monthly_outputs: "Monthly Outputs ($)",

    // Transactions
    transactions_title: "Transactions",
    btn_new_transaction: "+ New Transaction",
    form_create_transaction: "Create Transaction",
    form_edit_transaction: "Edit Transaction",
    label_description: "Description",
    ph_description: "Enter description",
    label_amount: "Amount ($)",
    label_is_debt: "Is Debt",
    label_user_id: "User ID",
    opt_select_user: "Select a user...",
    opt_failed_users: "Failed to load users",

    // Goals
    goals_title: "Goals",
    btn_new_goal: "+ New Goal",
    form_create_goal: "Create Goal",
    form_edit_goal: "Edit Goal",
    label_name: "Name",
    ph_goal_name: "Goal name",
    ph_goal_description: "Goal description",
    label_price: "Price ($)",
    label_pros: "Pros",
    ph_pros: "Advantages of this goal",
    label_cons: "Cons",
    ph_cons: "Disadvantages of this goal",
    label_deadline: "Deadline",

    // Buttons
    btn_create: "Create",
    btn_update: "Update",
    btn_cancel: "Cancel",
    btn_delete: "Delete",
    btn_change_user: "Change",
    btn_go_to_users: "Go to Users",

    // Table headers
    th_username: "Username",
    th_current_amount: "Current Amount",
    th_monthly_inputs: "Monthly Inputs",
    th_monthly_outputs: "Monthly Outputs",
    th_actions: "Actions",
    th_description: "Description",
    th_amount: "Amount",
    th_debt: "Debt",
    th_user_id: "User ID",
    th_created_at: "Created At",
    th_name: "Name",
    th_price: "Price",
    th_pros: "Pros",
    th_cons: "Cons",
    th_deadline: "Deadline",

    // Status / empty messages
    loading: "Loading...",
    no_users: "No users found. Create one to get started!",
    no_transactions: "No transactions found for this user.",
    no_goals: "No goals found for this user.",
    fail_users: "Failed to load users",
    fail_transactions: "Failed to load transactions",
    fail_goals: "Failed to load goals",

    // Toasts
    toast_user_created: "User created successfully",
    toast_user_updated: "User updated successfully",
    toast_user_deleted: "User deleted successfully",
    toast_transaction_created: "Transaction created successfully",
    toast_transaction_updated: "Transaction updated successfully",
    toast_transaction_deleted: "Transaction deleted successfully",
    toast_goal_created: "Goal created successfully",
    toast_goal_updated: "Goal updated successfully",
    toast_goal_deleted: "Goal deleted successfully",
    toast_select_user: "Please select a user",
    toast_user_selected: 'User "{name}" selected',

    // Confirm dialogs
    confirm_delete_user:
      'Are you sure you want to delete user "{name}" (ID: {id})?',
    confirm_delete_transaction:
      "Are you sure you want to delete transaction #{id}?",
    confirm_delete_goal:
      'Are you sure you want to delete goal "{name}" (ID: {id})?',

    // Badges
    badge_yes: "Yes",
    badge_no: "No",

    // Footer
    footer_text: "Take control of your finances",

    // Currency
    currency_symbol: "$",
    date_locale: "en-US",

    // Selected user
    no_user_selected: "No user selected",
    selected_user_label: "Selected:",
    hint_click_user:
      "Click on a user row to select it and view their transactions and goals.",
    prompt_select_user_transactions:
      "Select a user from the Users tab to view their transactions.",
    prompt_select_user_goals:
      "Select a user from the Users tab to view their goals.",
  },
  pt: {
    // Header
    subtitle: "Gerenciador Financeiro Pessoal",

    // Tabs
    tab_users: "Usuários",
    tab_transactions: "Transações",
    tab_goals: "Objetivos",

    // Users
    users_title: "Usuários",
    btn_new_user: "+ Novo Usuário",
    form_create_user: "Criar Usuário",
    form_edit_user: "Editar Usuário",
    label_username: "Nome de Usuário",
    ph_username: "Digite o nome de usuário",
    label_current_amount: "Saldo Atual (R$)",
    label_monthly_inputs: "Entradas Mensais (R$)",
    label_monthly_outputs: "Saídas Mensais (R$)",

    // Transactions
    transactions_title: "Transações",
    btn_new_transaction: "+ Nova Transação",
    form_create_transaction: "Criar Transação",
    form_edit_transaction: "Editar Transação",
    label_description: "Descrição",
    ph_description: "Digite a descrição",
    label_amount: "Valor (R$)",
    label_is_debt: "É Dívida",
    label_user_id: "ID do Usuário",
    opt_select_user: "Selecione um usuário...",
    opt_failed_users: "Falha ao carregar usuários",

    // Goals
    goals_title: "Objetivos",
    btn_new_goal: "+ Novo Objetivo",
    form_create_goal: "Criar Objetivo",
    form_edit_goal: "Editar Objetivo",
    label_name: "Nome",
    ph_goal_name: "Nome do objetivo",
    ph_goal_description: "Descrição do objetivo",
    label_price: "Preço (R$)",
    label_pros: "Prós",
    ph_pros: "Vantagens deste objetivo",
    label_cons: "Contras",
    ph_cons: "Desvantagens deste objetivo",
    label_deadline: "Prazo",

    // Buttons
    btn_create: "Criar",
    btn_update: "Atualizar",
    btn_cancel: "Cancelar",
    btn_delete: "Excluir",
    btn_change_user: "Trocar",
    btn_go_to_users: "Ir para Usuários",

    // Table headers
    th_username: "Nome de Usuário",
    th_current_amount: "Saldo Atual",
    th_monthly_inputs: "Entradas Mensais",
    th_monthly_outputs: "Saídas Mensais",
    th_actions: "Ações",
    th_description: "Descrição",
    th_amount: "Valor",
    th_debt: "Dívida",
    th_user_id: "ID do Usuário",
    th_created_at: "Criado Em",
    th_name: "Nome",
    th_price: "Preço",
    th_pros: "Prós",
    th_cons: "Contras",
    th_deadline: "Prazo",

    // Status / empty messages
    loading: "Carregando...",
    no_users: "Nenhum usuário encontrado. Crie um para começar!",
    no_transactions: "Nenhuma transação encontrada para este usuário.",
    no_goals: "Nenhum objetivo encontrado para este usuário.",
    fail_users: "Falha ao carregar usuários",
    fail_transactions: "Falha ao carregar transações",
    fail_goals: "Falha ao carregar objetivos",

    // Toasts
    toast_user_created: "Usuário criado com sucesso",
    toast_user_updated: "Usuário atualizado com sucesso",
    toast_user_deleted: "Usuário excluído com sucesso",
    toast_transaction_created: "Transação criada com sucesso",
    toast_transaction_updated: "Transação atualizada com sucesso",
    toast_transaction_deleted: "Transação excluída com sucesso",
    toast_goal_created: "Objetivo criado com sucesso",
    toast_goal_updated: "Objetivo atualizado com sucesso",
    toast_goal_deleted: "Objetivo excluído com sucesso",
    toast_select_user: "Por favor, selecione um usuário",
    toast_user_selected: 'Usuário "{name}" selecionado',

    // Confirm dialogs
    confirm_delete_user:
      'Tem certeza que deseja excluir o usuário "{name}" (ID: {id})?',
    confirm_delete_transaction:
      "Tem certeza que deseja excluir a transação #{id}?",
    confirm_delete_goal:
      'Tem certeza que deseja excluir o objetivo "{name}" (ID: {id})?',

    // Badges
    badge_yes: "Sim",
    badge_no: "Não",

    // Footer
    footer_text: "Assuma o controle das suas finanças",

    // Currency
    currency_symbol: "R$",
    date_locale: "pt-BR",

    // Selected user
    no_user_selected: "Nenhum usuário selecionado",
    selected_user_label: "Selecionado:",
    hint_click_user:
      "Clique em uma linha de usuário para selecioná-lo e ver suas transações e objetivos.",
    prompt_select_user_transactions:
      "Selecione um usuário na aba Usuários para ver suas transações.",
    prompt_select_user_goals:
      "Selecione um usuário na aba Usuários para ver seus objetivos.",
  },
};

let currentLang = localStorage.getItem(LANG_KEY) || "en";

/**
 * Returns the translated string for the given key.
 * Supports simple interpolation: t("key", { name: "x", id: 1 })
 */
function t(key, params) {
  const dict = translations[currentLang] || translations.en;
  let str = dict[key] || translations.en[key] || key;
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      str = str.replace(new RegExp(`\\{${k}\\}`, "g"), v);
    }
  }
  return str;
}

/**
 * Traverses the DOM and applies translations using data-i18n attributes.
 */
function applyLanguage() {
  document.documentElement.lang = currentLang === "pt" ? "pt-BR" : "en";

  // Translate text content
  document.querySelectorAll("[data-i18n]").forEach((el) => {
    const key = el.getAttribute("data-i18n");
    el.textContent = t(key);
  });

  // Translate placeholders
  document.querySelectorAll("[data-i18n-placeholder]").forEach((el) => {
    const key = el.getAttribute("data-i18n-placeholder");
    el.placeholder = t(key);
  });

  // Update flag button active states
  document
    .getElementById("lang-en")
    .classList.toggle("active", currentLang === "en");
  document
    .getElementById("lang-pt")
    .classList.toggle("active", currentLang === "pt");

  // Update selected user bar text
  updateSelectedUserBar();
}

function setLanguage(lang) {
  currentLang = lang;
  localStorage.setItem(LANG_KEY, lang);
  applyLanguage();
  refreshActiveTab();
}

function refreshActiveTab() {
  const activeTab = document.querySelector(".tab-btn.active");
  if (!activeTab) return;
  const target = activeTab.dataset.tab;
  if (target === "users") loadUsers();
  else if (target === "transactions") loadTransactions();
  else if (target === "goals") loadGoals();
}

// ===== Utility Helpers =====

function currencySymbol() {
  return t("currency_symbol");
}

function centsToDisplay(cents) {
  const value = cents / 100;
  return value.toLocaleString(t("date_locale"), {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
}

function dollarsToCents(dollars) {
  return Math.round(parseFloat(dollars) * 100);
}

function moneyClass(cents) {
  if (cents > 0) return "money positive";
  if (cents < 0) return "money negative";
  return "money";
}

function escapeHtml(str) {
  if (str == null) return "";
  const div = document.createElement("div");
  div.textContent = String(str);
  return div.innerHTML;
}

function formatDate(dateStr) {
  if (!dateStr) return "\u2014";
  try {
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return escapeHtml(dateStr);
    return d.toLocaleDateString(t("date_locale"), {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  } catch {
    return escapeHtml(dateStr);
  }
}

// ===== Toast Notifications =====

let toastTimeout = null;

function showToast(message, type = "success") {
  const toast = document.getElementById("toast");
  toast.textContent = message;
  toast.className = `toast ${type}`;
  toast.classList.remove("hidden");

  if (toastTimeout) clearTimeout(toastTimeout);
  toastTimeout = setTimeout(() => {
    toast.style.animation = "toastOut 0.3s ease forwards";
    setTimeout(() => {
      toast.classList.add("hidden");
      toast.style.animation = "";
    }, 300);
  }, 3000);
}

// ===== Confirm Dialog =====

let confirmCallback = null;

function showConfirm(message, onConfirm) {
  const overlay = document.getElementById("confirm-overlay");
  document.getElementById("confirm-message").textContent = message;
  overlay.classList.remove("hidden");
  confirmCallback = onConfirm;
}

document.getElementById("confirm-yes").addEventListener("click", () => {
  document.getElementById("confirm-overlay").classList.add("hidden");
  if (confirmCallback) confirmCallback();
  confirmCallback = null;
});

document.getElementById("confirm-no").addEventListener("click", () => {
  document.getElementById("confirm-overlay").classList.add("hidden");
  confirmCallback = null;
});

// ===== API Fetch Wrapper =====

async function apiFetch(url, options = {}) {
  try {
    const res = await fetch(url, {
      headers: { "Content-Type": "application/json" },
      ...options,
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(data.error || `HTTP ${res.status}`);
    }
    return data;
  } catch (err) {
    showToast(err.message || "Network error", "error");
    throw err;
  }
}

// ========================================================================
//  SELECTED USER STATE
// ========================================================================

let selectedUser = null; // { id, user_name }

function updateSelectedUserBar() {
  const bar = document.getElementById("selected-user-bar");
  const nameEl = document.getElementById("selected-user-name");

  if (selectedUser) {
    nameEl.textContent =
      t("selected_user_label") + " " + selectedUser.user_name;
    nameEl.removeAttribute("data-i18n");
    bar.classList.remove("hidden");
  } else {
    bar.classList.add("hidden");
  }
}

function selectUser(id, userName) {
  selectedUser = { id: id, user_name: userName };
  updateSelectedUserBar();
  highlightSelectedRow();
  showToast(t("toast_user_selected", { name: userName }), "info");
}

function deselectUser() {
  selectedUser = null;
  updateSelectedUserBar();
  highlightSelectedRow();
  updateTransactionsView();
  updateGoalsView();
}

function highlightSelectedRow() {
  document.querySelectorAll("#users-tbody tr").forEach((row) => {
    row.classList.remove("selected-row");
  });
  if (selectedUser) {
    const row = document.querySelector(
      `#users-tbody tr[data-user-id="${selectedUser.id}"]`,
    );
    if (row) row.classList.add("selected-row");
  }
}

function updateTransactionsView() {
  const noUser = document.getElementById("transactions-no-user");
  const content = document.getElementById("transactions-content");
  if (selectedUser) {
    noUser.classList.add("hidden");
    content.classList.remove("hidden");
  } else {
    noUser.classList.remove("hidden");
    content.classList.add("hidden");
  }
}

function updateGoalsView() {
  const noUser = document.getElementById("goals-no-user");
  const content = document.getElementById("goals-content");
  if (selectedUser) {
    noUser.classList.add("hidden");
    content.classList.remove("hidden");
  } else {
    noUser.classList.remove("hidden");
    content.classList.add("hidden");
  }
}

// Deselect button
document.getElementById("btn-deselect-user").addEventListener("click", () => {
  deselectUser();
  // Switch to users tab
  switchToTab("users");
});

// "Go to Users" buttons on the no-user prompts
document
  .getElementById("btn-go-users-from-transactions")
  .addEventListener("click", () => {
    switchToTab("users");
  });

document
  .getElementById("btn-go-users-from-goals")
  .addEventListener("click", () => {
    switchToTab("users");
  });

// ===== Tab Navigation =====

const tabButtons = document.querySelectorAll(".tab-btn");
const tabContents = document.querySelectorAll(".tab-content");

function switchToTab(tabName) {
  tabButtons.forEach((b) => b.classList.remove("active"));
  tabContents.forEach((c) => c.classList.remove("active"));

  const btn = document.querySelector(`.tab-btn[data-tab="${tabName}"]`);
  if (btn) btn.classList.add("active");
  document.getElementById(`tab-${tabName}`).classList.add("active");

  if (tabName === "users") loadUsers();
  else if (tabName === "transactions") loadTransactions();
  else if (tabName === "goals") loadGoals();
}

tabButtons.forEach((btn) => {
  btn.addEventListener("click", () => {
    switchToTab(btn.dataset.tab);
  });
});

// ========================================================================
//  USERS CRUD
// ========================================================================

const userFormContainer = document.getElementById("user-form-container");
const userForm = document.getElementById("user-form");
const userFormTitle = document.getElementById("user-form-title");
const userSubmitBtn = document.getElementById("user-submit-btn");
const usersTbody = document.getElementById("users-tbody");

let editingUserId = null;

document.getElementById("btn-add-user").addEventListener("click", () => {
  editingUserId = null;
  userFormTitle.textContent = t("form_create_user");
  userSubmitBtn.textContent = t("btn_create");
  userForm.reset();
  document.getElementById("user-id").value = "";
  userFormContainer.classList.remove("hidden");
});

document.getElementById("user-cancel-btn").addEventListener("click", () => {
  userFormContainer.classList.add("hidden");
  userForm.reset();
  editingUserId = null;
});

userForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const userName = document.getElementById("user-name").value.trim();
  const currentAmount = dollarsToCents(
    document.getElementById("user-current-amount").value,
  );
  const monthlyInputs = dollarsToCents(
    document.getElementById("user-monthly-inputs").value,
  );
  const monthlyOutputs = dollarsToCents(
    document.getElementById("user-monthly-outputs").value,
  );

  try {
    if (editingUserId) {
      const body = {};
      body.user_name = userName;
      body.current_amount = currentAmount;
      body.monthly_inputs = monthlyInputs;
      body.monthly_outputs = monthlyOutputs;

      await apiFetch(`${API.users}/${editingUserId}`, {
        method: "PATCH",
        body: JSON.stringify(body),
      });
      showToast(t("toast_user_updated"));

      // Update selectedUser name if editing the selected user
      if (selectedUser && selectedUser.id === editingUserId) {
        selectedUser.user_name = userName;
        updateSelectedUserBar();
      }
    } else {
      await apiFetch(API.users, {
        method: "POST",
        body: JSON.stringify({
          user_name: userName,
          current_amount: currentAmount,
          monthly_inputs: monthlyInputs,
          monthly_outputs: monthlyOutputs,
        }),
      });
      showToast(t("toast_user_created"));
    }

    userFormContainer.classList.add("hidden");
    userForm.reset();
    editingUserId = null;
    await loadUsers();
  } catch {
    // error already shown in toast
  }
});

async function loadUsers() {
  try {
    const users = await apiFetch(API.users);
    renderUsers(users);
  } catch {
    usersTbody.innerHTML = `<tr><td colspan="6" class="empty-msg">${escapeHtml(t("fail_users"))}</td></tr>`;
  }
}

function renderUsers(users) {
  if (!users || users.length === 0) {
    usersTbody.innerHTML = `<tr><td colspan="6" class="empty-msg">${escapeHtml(t("no_users"))}</td></tr>`;
    return;
  }

  const sym = currencySymbol();
  usersTbody.innerHTML = users
    .map(
      (u) => `
        <tr data-user-id="${u.id}" onclick="handleUserRowClick(${u.id}, '${escapeHtml(u.user_name).replace(/'/g, "\\'")}')"
            class="${selectedUser && selectedUser.id === u.id ? "selected-row" : ""}">
            <td>${u.id}</td>
            <td>${escapeHtml(u.user_name)}</td>
            <td><span class="${moneyClass(u.current_amount)}">${sym}${centsToDisplay(u.current_amount)}</span></td>
            <td><span class="${moneyClass(u.monthly_inputs)}">${sym}${centsToDisplay(u.monthly_inputs)}</span></td>
            <td><span class="money negative">${sym}${centsToDisplay(u.monthly_outputs)}</span></td>
            <td>
                <div class="actions-cell">
                    <button class="btn btn-icon edit" title="Edit" onclick="event.stopPropagation(); editUser(${u.id})">✏️</button>
                    <button class="btn btn-icon delete" title="Delete" onclick="event.stopPropagation(); deleteUser(${u.id}, '${escapeHtml(u.user_name).replace(/'/g, "\\'")}')">🗑️</button>
                </div>
            </td>
        </tr>
    `,
    )
    .join("");
}

window.handleUserRowClick = function (id, userName) {
  selectUser(id, userName);
};

window.editUser = async function (id) {
  try {
    const user = await apiFetch(`${API.users}/${id}`);
    editingUserId = id;
    userFormTitle.textContent = t("form_edit_user");
    userSubmitBtn.textContent = t("btn_update");
    document.getElementById("user-id").value = user.id;
    document.getElementById("user-name").value = user.user_name;
    document.getElementById("user-current-amount").value = (
      user.current_amount / 100
    ).toFixed(2);
    document.getElementById("user-monthly-inputs").value = (
      user.monthly_inputs / 100
    ).toFixed(2);
    document.getElementById("user-monthly-outputs").value = (
      user.monthly_outputs / 100
    ).toFixed(2);
    userFormContainer.classList.remove("hidden");
  } catch {
    // error already shown
  }
};

window.deleteUser = function (id, name) {
  showConfirm(t("confirm_delete_user", { name: name, id: id }), async () => {
    try {
      await apiFetch(`${API.users}/${id}`, { method: "DELETE" });
      showToast(t("toast_user_deleted"));
      // If the deleted user was selected, deselect
      if (selectedUser && selectedUser.id === id) {
        deselectUser();
      }
      await loadUsers();
    } catch {
      // error already shown
    }
  });
};

// ========================================================================
//  TRANSACTIONS CRUD
// ========================================================================

const transactionFormContainer = document.getElementById(
  "transaction-form-container",
);
const transactionForm = document.getElementById("transaction-form");
const transactionFormTitle = document.getElementById("transaction-form-title");
const transactionSubmitBtn = document.getElementById("transaction-submit-btn");
const transactionsTbody = document.getElementById("transactions-tbody");

let editingTransactionId = null;

document.getElementById("btn-add-transaction").addEventListener("click", () => {
  if (!selectedUser) {
    showToast(t("toast_select_user"), "error");
    return;
  }
  editingTransactionId = null;
  transactionFormTitle.textContent = t("form_create_transaction");
  transactionSubmitBtn.textContent = t("btn_create");
  transactionForm.reset();
  document.getElementById("transaction-id").value = "";
  transactionFormContainer.classList.remove("hidden");
});

document
  .getElementById("transaction-cancel-btn")
  .addEventListener("click", () => {
    transactionFormContainer.classList.add("hidden");
    transactionForm.reset();
    editingTransactionId = null;
  });

transactionForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const description = document
    .getElementById("transaction-description")
    .value.trim();
  const amount = dollarsToCents(
    document.getElementById("transaction-amount").value,
  );
  const isDebt = document.getElementById("transaction-is-debt").checked;

  try {
    if (editingTransactionId) {
      const body = {};
      body.description = description;
      body.amount = amount;
      body.is_debt = isDebt;

      await apiFetch(`${API.transactions}/${editingTransactionId}`, {
        method: "PATCH",
        body: JSON.stringify(body),
      });
      showToast(t("toast_transaction_updated"));
    } else {
      if (!selectedUser) {
        showToast(t("toast_select_user"), "error");
        return;
      }
      await apiFetch(API.transactions, {
        method: "POST",
        body: JSON.stringify({
          description: description,
          amount: amount,
          is_debt: isDebt,
          user_id: selectedUser.id,
        }),
      });
      showToast(t("toast_transaction_created"));
    }

    transactionFormContainer.classList.add("hidden");
    transactionForm.reset();
    editingTransactionId = null;
    await loadTransactions();
  } catch {
    // error already shown
  }
});

async function loadTransactions() {
  updateTransactionsView();
  if (!selectedUser) return;

  try {
    const transactions = await apiFetch(
      `${API.users}/transactions/${selectedUser.id}`,
    );
    renderTransactions(transactions);
  } catch {
    transactionsTbody.innerHTML = `<tr><td colspan="6" class="empty-msg">${escapeHtml(t("fail_transactions"))}</td></tr>`;
  }
}

function renderTransactions(transactions) {
  if (!transactions || transactions.length === 0) {
    transactionsTbody.innerHTML = `<tr><td colspan="6" class="empty-msg">${escapeHtml(t("no_transactions"))}</td></tr>`;
    return;
  }

  const sym = currencySymbol();
  const badgeYes = escapeHtml(t("badge_yes"));
  const badgeNo = escapeHtml(t("badge_no"));
  transactionsTbody.innerHTML = transactions
    .map(
      (txn) => `
        <tr>
            <td>${txn.id}</td>
            <td class="truncate" title="${escapeHtml(txn.description)}">${escapeHtml(txn.description) || "\u2014"}</td>
            <td><span class="${moneyClass(txn.amount)}">${sym}${centsToDisplay(txn.amount)}</span></td>
            <td><span class="badge ${txn.is_debt ? "badge-yes" : "badge-no"}">${txn.is_debt ? badgeYes : badgeNo}</span></td>
            <td class="date-cell">${formatDate(txn.created_at)}</td>
            <td>
                <div class="actions-cell">
                    <button class="btn btn-icon edit" title="Edit" onclick="editTransaction(${txn.id})">✏️</button>
                    <button class="btn btn-icon delete" title="Delete" onclick="deleteTransaction(${txn.id})">🗑️</button>
                </div>
            </td>
        </tr>
    `,
    )
    .join("");
}

window.editTransaction = async function (id) {
  try {
    const txn = await apiFetch(`${API.transactions}/${id}`);
    editingTransactionId = id;
    transactionFormTitle.textContent = t("form_edit_transaction");
    transactionSubmitBtn.textContent = t("btn_update");
    document.getElementById("transaction-id").value = txn.id;
    document.getElementById("transaction-description").value =
      txn.description || "";
    document.getElementById("transaction-amount").value = (
      txn.amount / 100
    ).toFixed(2);
    document.getElementById("transaction-is-debt").checked = txn.is_debt;
    transactionFormContainer.classList.remove("hidden");
  } catch {
    // error already shown
  }
};

window.deleteTransaction = function (id) {
  showConfirm(t("confirm_delete_transaction", { id: id }), async () => {
    try {
      await apiFetch(`${API.transactions}/${id}`, {
        method: "DELETE",
      });
      showToast(t("toast_transaction_deleted"));
      await loadTransactions();
    } catch {
      // error already shown
    }
  });
};

// ========================================================================
//  GOALS CRUD
// ========================================================================

const goalFormContainer = document.getElementById("goal-form-container");
const goalForm = document.getElementById("goal-form");
const goalFormTitle = document.getElementById("goal-form-title");
const goalSubmitBtn = document.getElementById("goal-submit-btn");
const goalsTbody = document.getElementById("goals-tbody");

let editingGoalId = null;

document.getElementById("btn-add-goal").addEventListener("click", () => {
  if (!selectedUser) {
    showToast(t("toast_select_user"), "error");
    return;
  }
  editingGoalId = null;
  goalFormTitle.textContent = t("form_create_goal");
  goalSubmitBtn.textContent = t("btn_create");
  goalForm.reset();
  document.getElementById("goal-id").value = "";
  goalFormContainer.classList.remove("hidden");
});

document.getElementById("goal-cancel-btn").addEventListener("click", () => {
  goalFormContainer.classList.add("hidden");
  goalForm.reset();
  editingGoalId = null;
});

goalForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const name = document.getElementById("goal-name").value.trim();
  const description = document.getElementById("goal-description").value.trim();
  const price = dollarsToCents(document.getElementById("goal-price").value);
  const pros = document.getElementById("goal-pros").value.trim();
  const cons = document.getElementById("goal-cons").value.trim();
  const deadline = document.getElementById("goal-deadline").value;

  try {
    if (editingGoalId) {
      const body = {};
      body.name = name;
      if (description) body.description = description;
      body.price = price;
      if (pros) body.pros = pros;
      if (cons) body.cons = cons;
      if (deadline) body.deadline = deadline;

      await apiFetch(`${API.goals}/${editingGoalId}`, {
        method: "PATCH",
        body: JSON.stringify(body),
      });
      showToast(t("toast_goal_updated"));
    } else {
      if (!selectedUser) {
        showToast(t("toast_select_user"), "error");
        return;
      }
      await apiFetch(API.goals, {
        method: "POST",
        body: JSON.stringify({
          name: name,
          description: description,
          price: price,
          pros: pros,
          cons: cons,
          user_id: selectedUser.id,
          deadline: deadline,
        }),
      });
      showToast(t("toast_goal_created"));
    }

    goalFormContainer.classList.add("hidden");
    goalForm.reset();
    editingGoalId = null;
    await loadGoals();
  } catch {
    // error already shown
  }
});

async function loadGoals() {
  updateGoalsView();
  if (!selectedUser) return;

  try {
    const goals = await apiFetch(`${API.users}/goals/${selectedUser.id}`);
    renderGoals(goals);
  } catch {
    goalsTbody.innerHTML = `<tr><td colspan="9" class="empty-msg">${escapeHtml(t("fail_goals"))}</td></tr>`;
  }
}

function renderGoals(goals) {
  if (!goals || goals.length === 0) {
    goalsTbody.innerHTML = `<tr><td colspan="9" class="empty-msg">${escapeHtml(t("no_goals"))}</td></tr>`;
    return;
  }

  const sym = currencySymbol();
  goalsTbody.innerHTML = goals
    .map(
      (g) => `
        <tr>
            <td>${g.id}</td>
            <td>${escapeHtml(g.name)}</td>
            <td class="truncate" title="${escapeHtml(g.description)}">${escapeHtml(g.description) || "\u2014"}</td>
            <td><span class="${moneyClass(g.price)}">${sym}${centsToDisplay(g.price)}</span></td>
            <td class="truncate" title="${escapeHtml(g.pros)}">${escapeHtml(g.pros) || "\u2014"}</td>
            <td class="truncate" title="${escapeHtml(g.cons)}">${escapeHtml(g.cons) || "\u2014"}</td>
            <td class="date-cell">${formatDate(g.deadline)}</td>
            <td class="date-cell">${formatDate(g.created_at)}</td>
            <td>
                <div class="actions-cell">
                    <button class="btn btn-icon edit" title="Edit" onclick="editGoal(${g.id})">✏️</button>
                    <button class="btn btn-icon delete" title="Delete" onclick="deleteGoal(${g.id}, '${escapeHtml(g.name).replace(/'/g, "\\'")}')">🗑️</button>
                </div>
            </td>
        </tr>
    `,
    )
    .join("");
}

window.editGoal = async function (id) {
  try {
    const g = await apiFetch(`${API.goals}/${id}`);
    editingGoalId = id;
    goalFormTitle.textContent = t("form_edit_goal");
    goalSubmitBtn.textContent = t("btn_update");
    document.getElementById("goal-id").value = g.id;
    document.getElementById("goal-name").value = g.name || "";
    document.getElementById("goal-description").value = g.description || "";
    document.getElementById("goal-price").value = (g.price / 100).toFixed(2);
    document.getElementById("goal-pros").value = g.pros || "";
    document.getElementById("goal-cons").value = g.cons || "";
    // Set deadline
    if (g.deadline) {
      try {
        const d = new Date(g.deadline);
        if (!isNaN(d.getTime())) {
          const yyyy = d.getFullYear();
          const mm = String(d.getMonth() + 1).padStart(2, "0");
          const dd = String(d.getDate()).padStart(2, "0");
          document.getElementById("goal-deadline").value =
            `${yyyy}-${mm}-${dd}`;
        } else {
          document.getElementById("goal-deadline").value = g.deadline;
        }
      } catch {
        document.getElementById("goal-deadline").value = g.deadline;
      }
    }
    goalFormContainer.classList.remove("hidden");
  } catch {
    // error already shown
  }
};

window.deleteGoal = function (id, name) {
  showConfirm(t("confirm_delete_goal", { name: name, id: id }), async () => {
    try {
      await apiFetch(`${API.goals}/${id}`, { method: "DELETE" });
      showToast(t("toast_goal_deleted"));
      await loadGoals();
    } catch {
      // error already shown
    }
  });
};

// ========================================================================
//  LANGUAGE SWITCHER EVENT LISTENERS
// ========================================================================

document.getElementById("lang-en").addEventListener("click", () => {
  setLanguage("en");
});

document.getElementById("lang-pt").addEventListener("click", () => {
  setLanguage("pt");
});

// ========================================================================
//  INITIAL LOAD
// ========================================================================

document.addEventListener("DOMContentLoaded", () => {
  applyLanguage();
  updateTransactionsView();
  updateGoalsView();
  loadUsers();
});
