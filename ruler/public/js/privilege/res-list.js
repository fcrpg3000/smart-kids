/* 
 * Copyright (C) 2012-2013 king4go authors All rights reserved.
 *
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements. See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 *           http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

(function($) {
  var jDel, jDelModal, jDelModalBody, delId;
  function initialize() {
    jDel = $('#btn_del_res');
    jDelModal = $('#modal_del_res');
    jDelModalBody = $('div.modal-body', jDelModal);

    jDel.click(function() {
      $.post('/privilege/a/del_res',
        {'id': delId}, 
        function(data) {
          alert(data.message);
          if (data.code === 1) {
            location.reload();
          }
        }, 'json');
      return false;
    });
  }
  function preDelete(id, name) {
    delId = id;
    jDelModalBody.html([
      '<p>你确定要删除<strong>“', name, '”</strong>这个资源吗？</p>',
      '<p><strong>警告！</strong>资源删除后，所有关联的角色将不能访问该资源！</p>'
      ].join(''));
    jDelModal.modal('show');
    return false;
  }
  $(function() {
    initialize();
  });
  window.resPreDelete = preDelete;
})(jQuery);