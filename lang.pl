sealed interface Body {}
record Text(body: String) implement Body {}
record Markdown(body: MarkdownString) implement Body {}
record Html(body: HtmlString) implement Body {}

class Mailer {
  #"Split method (multipart method)"
  fun sendTo(receiver: Email) 
      subject(subject: String)
      body(body: Body)
      copy(receivers: Email...)
  {
  }

  #MyAnnotation [ property = "value" ]
  fun close() {
  }
}

#"Leading single line documentation"
Mailer().sendTo(Email("john@mail.com"))
         subject("My subject")
         body(Text("A simple message body."))
         copy(Email("alice@mail.com"), Email("bob@mail.com"))
        .close();
#^"Trailing single line documentation"

def:block UnitBlock<T> {
  execute(t: T);
}

def:block VoidBlock<T> {
  execute();
}

def:record Person(name: String, dob: LocalDate);

def:fn calculatePrice(products: List<Product>, int quantity) {}

def:decorator User(first: Name, last: Name);

def-block VoidBlock<T> {
  execute();
}

def-block UnitBlock<T> {
  execute(t: T);
}

def-record Person(name: String, dob: LocalDate);

def-fn calculatePrice(products: List<Product>, int quantity) {}

def-decorator User(first: Name, last: Name);

sealed interface Boolean {
    fun ifTrue() 
        block then: VoidBlock<Boolean>;
    
    fun ifFalse() 
        block then: VoidBlock<Boolean>;
}

singleton True implements Boolean {
    override fun ifTrue() block then: VoidBlock<Boolean>
    { then.execute();
    }
    
    override fun ifFalse() block then: VoidBlock<Boolean> {
    }
}

singleton False implements Boolean {
    override fun ifTrue() block then: VoidBlock<Boolean> {
    }
    
    override fun ifFalse() block then: VoidBlock<Boolean>
    { then.execute();
    }
}

#!"Free (or unbound) single line documentation"

sealed interface Optional<T> {
	#"syntax 1"
	fun ifPresent() block then: UnitBlock<T>
	    otherwise() block elseBlock: VoidBlock<T>;
	
	#"syntax 2"
	fun ifPresent() 
	    block [ then: UnitBlock<T>, otherwise: VoidBlock<T> ];
	
	#"syntax 3"
	fun ifPresent()
	    block then: UnitBlock<T>
	    block otherwise: VoidBlock<T>;
	    
	#"syntax 4"
	block fun() 
	      ifPresent: UnitBlock<T>
	      otherwise: VoidBlock<T>;
}

#"syntax 1"
Optional.none()
  .ifPresent() { play();	}
   otherwise() { introduceYourself(); }

#"syntax 2"
Optional.none()
	.ifPresent() { play(); }
	 otherwise { introduceYourself(); }

#"syntax 3"
Optional.none()
	.ifPresent() { play(); }
	 otherwise { introduceYourself(); }

#"syntax 4"
Optional.none()
	.ifPresent { play(); }
	 otherwise { introduceYourself(); }

fun findById(id: String): User? #^"method that takes one String parameter and returns a User or null"
    this [name, age] #^"method that can be called in a constructor"
    let [ db: Repository, req: HttpReques, res: HttpResponse ] #^"Scoped variables"
    block [ then: VoidBlock<Boolean>, else: VoidBlock<Boolean> ] #^"Blocks"
    raises [ DbException, IoException ] #^"Exceptions";

fun findById(id: String): User? #^"method that takes one String parameter and returns a User or null"
    this [name, age] #^"method that can be called in a constructor"
    let dp: Repository #^"Scope variable"
    let req: HttpRequest #^"Scope variable"
    let res: HttpResponse #^"Scope variable"
    block then: VoidBlock<Boolean> #^"Block"
    block else: VoidBlock<Boolean> #^"Block"
    raises [ DbException, IoException ] #^"Exceptions"

#"Split method (multipart method)"
fun sendTo(receiver: Email) 
    subject(subject: String)
    body(body: String)
{ let mail = Email(receiver)
; mail.subject(subject)
; mail.body(body)
; mail.send()
}