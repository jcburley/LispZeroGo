(defglobal 'onelet
  '((me formlist env)
    (eval
     (eval (car (cdr formlist)) env)
     (cons
      (cons
       (eval (car (car (car formlist))) env)
       (eval (car (cdr (car (car formlist)))) env))
      env))
    ))

(defglobal 'null
  '((me formlist env)
    (onelet (('arg (eval (car formlist) env)))
	    '(cond ((atom arg) (eq arg '()))
		   ('t '())))))

; Returns a list consisting of the (evaluated) arguments (supports any number of args)
(defglobal 'list
  '((me formlist env)
    (cond ((null formlist) '())
	  ('t (cons (eval (car formlist) env)
		    (apply me me (cdr formlist) env))))))

; (applied params argsname envname body)
; TODO: Pull into caller(s)
(defglobal 'applied
  '((me formlist env)
    (onelet (('params (eval (car formlist) env)))
	    '(onelet (('argsname (eval (car (cdr formlist)) env)))
		     '(onelet (('envname (eval (car (cdr (cdr formlist))) env)))
			      '(onelet (('body (eval (car (cdr (cdr (cdr formlist)))) env)))
				       '(cond ((null params) body)
					      ('t (list 'onelet
							(list (list (list 'quote
									  (car params))
								    (list 'eval
									  (list 'car
										argsname)
									  envname)))
							(list 'quote
							      (applied (cdr params)
								       (list 'cdr
									     argsname)
								       envname
								       body)))))))))))

; (defun name params body)
(defglobal 'defun
  '((me formlist env)
    (defglobal
      (car formlist)
      (list '(me formlist env)
	    (applied (car (cdr formlist))
		     'formlist
		     'env
		     (car (cdr (cdr formlist))))))))

(.symbol_dump)
(.load)
